package main

import (
	"fmt"

	"xform/resize"
	"xform/takeout"
)

func main() {

	var (
		src   = "/home/trimble/takeout01"
		dst   = "/home/trimble/takeout01/resizedddd"
		sizes = resize.Sizes{
			{Name: "large", Scale: 4},
			{Name: "thumb", Scale: 16},
			{Name: "thumb-gs", Scale: 16, Gs: true},
		}
	)

	// Todo: Find is overkill here and only uses photo.Path which is only used here, plz refactor plz
	//       also use Name, maybe just leave it for now??

	photos, err := takeout.Find(src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("found %d photos\n", len(photos))

	// Todo: maybe want tmpdir??
	//tmp, err := os.MkdirTemp("/tmp", fmt.Sprintf("photos-%d-", scale))

	err = sizes.BulkResize(dst, photos)
	if err != nil {
		panic(err)
	}
}
