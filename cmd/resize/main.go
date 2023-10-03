package main

import (
	"fmt"

	"xform/resize"
	"xform/takeout"
)

func main() {

	var (
		src   = "/home/trimble/takeout01"
		dst   = "/home/trimble/takeout01/resized"
		sizes = resize.Sizes{
			//{Name: "large", Scale: 4},
			//{Name: "thumb", Scale: 16},
			{Name: "thumb-gs", Scale: 16, Gs: true},
		}
	)

	// Todo: Find is overkill here and only uses photo.Path which is only used here, plz refactor plz
	photos, err := takeout.Find(src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("found %d photos\n", len(photos))

	//	fmt.Printf(">>> %#v\n", photos[0])

	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-", scale))
	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-gs-", scale))
	//if err != nil {
	//panic(err)
	//}
	//tmpDir := "no save"
	err = sizes.BulkResize(dst, photos)
	if err != nil {
		panic(err)
	}

	//resize.Bulk(tmp, &photos, scale)
	//fmt.Printf(">>> %#v\n", photos[0])
}
