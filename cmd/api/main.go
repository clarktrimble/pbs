package main

import (
	"context"
	"sync"
	"xform/bolt"
	"xform/chi"
	"xform/photosvc"

	"github.com/clarktrimble/delish"
	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/delish/graceful"
	"github.com/clarktrimble/delish/mid"
	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
)

const (
	cfgPrefix string = "pbapi"
)

var (
	version string
	wg      sync.WaitGroup
)

type Config struct {
	Version string         `json:"version" ignored:"true"`
	Bolt    *bolt.Config   `json:"bolt"`
	Server  *delish.Config `json:"server"`
}

func main() {

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)
	// Todo: fix delish/graceful WithTimeout is bug!!!?

	// create logger and initialize graceful

	lgr := &minlog.MinLog{} // Todo: sabot for trunc
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	ctx = graceful.Initialize(ctx, &wg, 6*cfg.Server.Timeout, lgr)

	// create router/handler, and server

	rtr := chi.New()

	handler := mid.LogResponse(lgr, rtr)
	handler = mid.LogRequest(lgr, hondo.Rand, handler)
	handler = mid.ReplaceCtx(ctx, handler)

	svr := cfg.Server.New(handler, lgr)

	// setup photo service layer

	repo, err := cfg.Bolt.New()
	launch.Check(ctx, lgr, err)
	defer repo.Close()

	photoSvc := &photosvc.PhotoSvc{
		Server: svr,
		Repo:   repo,
	}

	// register routes

	photoSvc.Register(rtr)
	_ = rtr.Set("GET", "/config", svr.ObjHandler("config", cfg))

	// delicious!

	svr.Start(ctx, &wg)
	graceful.Wait(ctx)
}
