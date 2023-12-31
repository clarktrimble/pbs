package main

import (
	"context"
	"os"
	"pbs/bolt"
	"pbs/chi"
	"pbs/photosvc"
	"sync"

	"github.com/clarktrimble/delish"
	"github.com/clarktrimble/delish/graceful"
	"github.com/clarktrimble/delish/mid"
	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
	"github.com/clarktrimble/sabot"
)

const (
	cfgPrefix string = "pb"
)

var (
	version string
	wg      sync.WaitGroup
)

type Config struct {
	Version  string         `json:"version" ignored:"true"`
	Truncate int            `json:"truncate" desc:"truncate log fields beyond length"`
	Bolt     *bolt.Config   `json:"bolt"`
	Server   *delish.Config `json:"server"`
}

func main() {

	// load config, setup logger

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	lgr := &sabot.Sabot{Writer: os.Stdout, MaxLen: cfg.Truncate}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	ctx = graceful.Initialize(ctx, &wg, lgr)

	// create router/handler, and server

	rtr := chi.New()

	handler := mid.LogResponse(lgr, rtr)
	handler = mid.LogRequest(lgr, hondo.Rand, handler)
	handler = mid.ReplaceCtx(ctx, handler)

	svr := cfg.Server.New(handler, lgr)

	// setup service layer and register routes

	repo, err := cfg.Bolt.New()
	launch.Check(ctx, lgr, err)
	defer repo.Close()

	photoSvc := &photosvc.PhotoSvc{
		Logger: lgr,
		Repo:   repo,
	}

	photoSvc.Register(rtr)
	rtr.Set("GET", "/config", svr.ObjHandler("config", cfg))

	// delicious!

	svr.Start(ctx, &wg)
	graceful.Wait(ctx)
}
