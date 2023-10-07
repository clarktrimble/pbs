package main

import (
	"context"
	"fmt"
	"os"

	"xform/moresvc"
	"xform/resize"
	"xform/takeout"

	"github.com/clarktrimble/delish/examples/api/minlog"
	"github.com/clarktrimble/giant"
	"github.com/clarktrimble/giant/logrt"
	"github.com/clarktrimble/giant/statusrt"
	"github.com/clarktrimble/launch"
)

type Config struct {
	Version     string        `json:"version" ignored:"true"`
	TakeoutPath string        `json:"takeout_path" desc:"path to takeout jpg's and json's" required:"true"`
	ResizedPath string        `json:"resized_path" desc:"path to resized png's" required:"true"`
	Filter      string        `json:"filter_regex" desc:"regex selecting matches" required:"true"`
	ApiClient   *giant.Config `json:"api_client"`
	DryRun      bool          `json:"dry_run" desc:"dig up metadata, but don't post"`
}

func main() {

	const (
		cfgPrefix string = "pbl"
	)
	var (
		version string
		sizes   = resize.Sizes{
			{Name: "large", Scale: 4},
			{Name: "thumb", Scale: 16},
			{Name: "thumb-gs", Scale: 16, Gs: true},
		}
	)

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	ctx := context.Background()
	lgr := &minlog.MinLog{}

	photos, err := takeout.FromFiles(cfg.TakeoutPath, cfg.Filter)
	if err != nil {
		fmt.Printf("error: %+v\n\n", err)
		os.Exit(1)
		// Todo: use lh Check if? logger shows
	}

	err = resize.AddResize(photos, cfg.ResizedPath, sizes)
	launch.Check(ctx, lgr, err)

	fmt.Printf("found %d photos\n", len(photos))

	if cfg.DryRun {
		fmt.Printf("%s\n", photos)
		return
	}

	fmt.Printf("posting to api\n")

	client := cfg.ApiClient.New()
	client.Use(&statusrt.StatusRt{})
	client.Use(&logrt.LogRt{Logger: lgr})

	photoSvc := &moresvc.Svc{Client: client}
	err = photoSvc.PostPhotos(ctx, photos)
	if err != nil {
		lgr.Error(ctx, "failed to post data", err)
	}
}
