package main

import (
	"fmt"

	"xform/resize"
	"xform/takeout"
)

func main() {

	sizes := resize.Sizes{
		{Scale: 4},
		{Scale: 16},
	}

	//photos, err := takeout.Find("/Users/trimble/takeout")
	photos, err := takeout.Find("/home/trimble/takeout01")
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
	tmpDir := "no save"
	err = sizes.Bulk(tmpDir, photos)
	if err != nil {
		panic(err)
	}

	//resize.Bulk(tmp, &photos, scale)

	fmt.Printf(">>> %#v\n", photos[0])
}
