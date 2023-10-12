package main

import (
	"context"
	"fmt"
	"os"

	"xform/resize"
	"xform/takeout"

	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
	"github.com/clarktrimble/sabot"
)

const (
	cfgPrefix string = "pbr" // Todo: goes in pkg?? (esp version?)
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
	Filter      string `json:"filter_regex" desc:"regex selecting matches" required:"true"`
	DryRun      bool   `json:"dry_run" desc:"dig up metadata, but don't post"`
}

func main() {

	// load config and setup logger

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	lgr := &sabot.Sabot{Writer: os.Stdout, MaxLen: cfg.Truncate}
	ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

	// Todo: Find is overkill here and only uses photo.Path which is only used here, plz refactor plz
	//       also use Name, maybe just leave it for now??
	//       also regex is hiding in there!!

	photos, err := takeout.FromFiles(cfg.TakeoutPath, cfg.Filter)
	launch.Check(ctx, lgr, err)

	if cfg.DryRun {
		fmt.Printf("%s\n", photos)
		return
	}

	fmt.Printf("found %d photos\n", len(photos))

	// Todo: maybe want tmpdir??
	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-", scale))

	err = sizes.BulkResize(cfg.ResizedPath, photos)
	launch.Check(ctx, lgr, err)
}
