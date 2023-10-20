package main

import (
	"context"
	"os"

	"github.com/clarktrimble/giant"
	"github.com/clarktrimble/giant/logrt"
	"github.com/clarktrimble/giant/statusrt"
	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
	"github.com/clarktrimble/sabot"

	"pbs/clientsvc"
	"pbs/resize"
	"pbs/takeout"
)

const (
	cfgPrefix string = "pb"
)

var (
	version string
	sizes   = resize.Sizes{
		{Name: "large", Scale: 4},
		{Name: "thumb", Scale: 16},
		{Name: "thumb-gs", Scale: 16, Gs: true},
	}
)

type Config struct {
	Version     string        `json:"version" ignored:"true"`
	Truncate    int           `json:"truncate" desc:"truncate log fields beyond length"`
	TakeoutPath string        `json:"takeout_path" desc:"path to takeout jpg's and json's" required:"true"`
	ResizedPath string        `json:"resized_path" desc:"path to resized png's" required:"true"`
	Match       string        `json:"match_regex" desc:"regex selecting matches" required:"true"`
	ApiClient   *giant.Config `json:"api_client"`
	DryRun      bool          `json:"dry_run" desc:"dig up metadata, but don't post"`
}

func main() {

	// load config and setup logger

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	lgr := &sabot.Sabot{Writer: os.Stdout, MaxLen: cfg.Truncate}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	// scan takeout and resized folders

	tos, err := takeout.ScanTakeout(cfg.TakeoutPath, cfg.Match)
	launch.Check(ctx, lgr, err)

	photos, err := tos.Photos()
	launch.Check(ctx, lgr, err)

	err = sizes.AddImages(photos, cfg.ResizedPath)
	launch.Check(ctx, lgr, err)

	// ship it!

	lgr.Info(ctx, "posting photos", "count", len(photos))
	if cfg.DryRun {
		lgr.Info(ctx, "just kidding", "dry_run", true)
		return
	}

	client := cfg.ApiClient.New()
	client.Use(&statusrt.StatusRt{})
	client.Use(&logrt.LogRt{Logger: lgr})

	photoSvc := &clientsvc.Svc{Client: client}
	err = photoSvc.PostPhotos(ctx, photos)
	if err != nil {
		lgr.Error(ctx, "failed to post data", err)
	}
}
