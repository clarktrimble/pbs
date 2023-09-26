package main

import (
	"context"
	"xform/moresvc"
	"xform/takeout"

	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/giant"
	"github.com/clarktrimble/giant/logrt"
	"github.com/clarktrimble/giant/statusrt"
)

var (
	//lat     = 58.38
	//lon     = 25.73
	baseUri = "http://localhost:8088"
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
	//lgr := &minLog{}
	lgr := &minlog.MinLog{}

	photos, err := takeout.Find("/home/trimble/takeout01")
	if err != nil {
		panic(err)
	}

	client := cfg.Client.New()
	client.Use(&statusrt.StatusRt{})
	client.Use(&logrt.LogRt{Logger: lgr})

	photoSvc := &moresvc.Svc{Client: client}
	//hourly, err := weatherSvc.GetHourly(ctx, lat, lon)
	err = photoSvc.PostPhotos(ctx, photos)
	if err != nil {
		//lgr.Error(ctx, "failed to get forcast data", err)
		//os.Exit(1)
		panic(err)
	}

	//hourly.Print()
}
