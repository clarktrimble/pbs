package main

import (
	"context"
	"sync"
	"time"
	"xform/bolt"
	"xform/photosvc"

	"github.com/clarktrimble/delish"
	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/delish/examples/api/minroute"
	"github.com/clarktrimble/delish/graceful"
	"github.com/clarktrimble/delish/mid"
	"github.com/clarktrimble/hondo"
)

var (
	version string
	wg      sync.WaitGroup
)

type Config struct {
	Version string         `json:"version" ignored:"true"`
	Server  *delish.Config `json:"server"`
}

func main() {

	// Todo: load with envconfig helper

	cfg := &Config{
		Version: version,
		Server: &delish.Config{
			Port:    8088,
			Timeout: 999 * time.Minute,
			// Todo: fix delish/graceful WithTimeout is bug!!!?
		},
	}

	// create logger and initialize graceful

	lgr := &minlog.MinLog{}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	ctx = graceful.Initialize(ctx, &wg, 6*cfg.Server.Timeout, lgr)

	// create router/handler, and server

	rtr := minroute.New(lgr)

	handler := mid.LogResponse(lgr, rtr)
	handler = mid.LogRequest(lgr, hondo.Rand, handler)
	handler = mid.ReplaceCtx(ctx, handler)

	svr := cfg.Server.New(handler, lgr)

	// setup photo service layer

	blt, err := bolt.New("photo.db", "photo")
	if err != nil {
		panic(err)
	}
	defer blt.Close()

	photoSvc := &photosvc.PhotoSvc{
		Server: svr,
		Repo:   blt,
	}

	// register routes

	photoSvc.Register(rtr)
	rtr.Set("GET", "/config", svr.ObjHandler("config", cfg))

	// delicious!

	svr.Start(ctx, &wg)
	graceful.Wait(ctx)

	// Todo: service shutdown with errors or sommat?  noooooo
}
