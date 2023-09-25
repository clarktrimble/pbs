package main

import (
	"context"
	"sync"
	"time"
	"xform/photosvc"

	"github.com/clarktrimble/delish"
	"github.com/clarktrimble/hondo"

	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/delish/examples/api/minroute"
	"github.com/clarktrimble/delish/graceful"
	"github.com/clarktrimble/delish/mid"
	//"github.com/clarktrimble/delish/cmd/api/demosvc"
	//"github.com/clarktrimble/delish/cmd/api/minlog"
	//"github.com/clarktrimble/delish/cmd/api/minroute"
)

var (
	version string
	wg      sync.WaitGroup
)

type Config struct {
	Version string         `json:"version" ignored:"true"`
	Server  *delish.Config `json:"server"`
}

// Todo: demo another goroutine

func main() {

	// usually load config with envconfig, but literal for demo

	cfg := &Config{
		Version: version,
		Server: &delish.Config{
			Port:    8088,
			Timeout: 10 * time.Second,
		},
	}

	// create logger and initialize graceful

	//lgr := &sabot.Sabot{
	//Writer: os.Stdout,
	//}
	lgr := &minlog.MinLog{}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	ctx = graceful.Initialize(ctx, &wg, 6*cfg.Server.Timeout, lgr)

	// create router/handler, and server

	rtr := minroute.New(lgr)

	handler := mid.LogResponse(lgr, rtr)
	handler = mid.LogRequest(lgr, hondo.Rand, handler)
	handler = mid.ReplaceCtx(ctx, handler)

	svr := cfg.Server.New(handler, lgr)

	// register route directly

	rtr.Set("GET", "/config", svr.ObjHandler("config", cfg))

	// or via service layer

	photosvc.NewAndReg(svr, rtr)

	// delicious!

	svr.Start(ctx, &wg)
	graceful.Wait(ctx)
}
