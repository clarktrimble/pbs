package main

import (
	"context"
	"os"

	"pbs/resize"
	"pbs/takeout"

	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
	"github.com/clarktrimble/sabot"
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
	Version     string `json:"version" ignored:"true"`
	Truncate    int    `json:"truncate" desc:"truncate log fields beyond length"`
	TakeoutPath string `json:"takeout_path" desc:"path to takeout jpg's and json's" required:"true"`
	ResizedPath string `json:"resized_path" desc:"path to resized png's" required:"true"`
	Match       string `json:"match_regex" desc:"regex selecting matches" required:"true"`
	DryRun      bool   `json:"dry_run" desc:"dig up metadata, but don't post"`
}

func main() {

	// load config and setup logger

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	lgr := &sabot.Sabot{Writer: os.Stdout, MaxLen: cfg.Truncate}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	// scan takeout folder

	tos, err := takeout.ScanTakeout(cfg.TakeoutPath, cfg.Match)
	launch.Check(ctx, lgr, err)

	photos := tos.PhotoFiles()

	// and resize!

	lgr.Info(ctx, "resizing photos", "count", len(photos))
	if cfg.DryRun {
		lgr.Info(ctx, "just kidding", "dry_run", true)
		return
	}

	// maybe want tmpdir? ala:
	// tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-", scale))

	err = sizes.ResizePhotos(cfg.ResizedPath, photos)
	launch.Check(ctx, lgr, err)
}
