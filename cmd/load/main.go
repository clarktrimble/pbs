package main

import (
	"context"
	"fmt"

	"xform/resize"
	"xform/takeouttoo"

	"github.com/clarktrimble/launch"
)

type Config struct {
	Version     string `json:"version" ignored:"true"`
	TakeoutPath string `json:"takeout_path" desc:"path to takeout jpg's and json's" required:"true"`
	ResizedPath string `json:"resized_path" desc:"path to resized png's" required:"true"`
}

func main() {

	const (
		cfgPrefix string = "pbl"
	)
	var (
		version string
		//src     = "/home/trimble/takeout01"
		//dst     = "/home/trimble/takeout01/resized"
		sizes = resize.Sizes{
			{Name: "large", Scale: 4},
			{Name: "thumb", Scale: 16},
			{Name: "thumb-gs", Scale: 16, Gs: true},
		}
	)

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix)

	//photos, err := takeouttoo.Find("/home/trimble/takeout01")
	//func FromFiles(jsonPath, resizePath string, sizes resize.Sizes) (photos entity.Photos, err error) {

	fmt.Printf(">>> %#v\n", cfg)
	photos, err := takeouttoo.FromFiles(cfg.TakeoutPath)
	if err != nil {
		launch.Check(context.Background(), nil, err)
	}
	return

	err = resize.AddResize(photos, cfg.ResizedPath, sizes)

	fmt.Printf("found %d photos\n", len(photos))
	//	fmt.Printf(">>> %#v\n", photos[0])

	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-", scale))
	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-gs-", scale))
	//if err != nil {
	//panic(err)
	//}
	//tmpDir := "no save"
	//err = sizes.BulkResize(tmpDir, photos)
	//if err != nil {
	//panic(err)
	//}

	//resize.Bulk(tmp, &photos, scale)

	fmt.Printf(">>> %#v\n", photos[0])
}
