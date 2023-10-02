package main

import (
	"context"
	"fmt"
	"xform/moresvc"
	_ "xform/resize" // Todo: wahh?
	"xform/takeout"

	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/giant"
	"github.com/clarktrimble/giant/logrt"
	"github.com/clarktrimble/giant/statusrt"
)

var (
	baseUri = "http://localhost:8088"
	path    = "/home/trimble/takeout01"
	dryRun  = false
)

type Config struct {
	Client *giant.Config
}

func main() {

	cfg := &Config{
		Client: &giant.Config{
			BaseUri: baseUri,
		},
	}

	ctx := context.Background()
	lgr := &minlog.MinLog{}

	photos, err := takeout.Find(path)
	if err != nil {
		panic(err)
	}

	if dryRun {
		fmt.Printf("\nfound:\n\n%s\n\n", photos)
		return
	}
	fmt.Printf("\nposting %d photos\n\n", len(photos))

	client := cfg.Client.New()
	client.Use(&statusrt.StatusRt{})
	client.Use(&logrt.LogRt{Logger: lgr})

	photoSvc := &moresvc.Svc{Client: client}
	err = photoSvc.PostPhotos(ctx, photos)
	if err != nil {
		//lgr.Error(ctx, "failed to get forcast data", err)
		//os.Exit(1)
		panic(err)
	}
}
