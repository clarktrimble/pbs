package main

import (
	"fmt"

	"xform/resize"
	"xform/takeouttoo"
)

func main() {

	//sizes := resize.Sizes{
	//{Scale: 4},
	//{Scale: 16},
	//}
	var (
		src   = "/home/trimble/takeout01"
		dst   = "/home/trimble/takeout01/resized"
		sizes = resize.Sizes{
			{Name: "large", Scale: 4},
			{Name: "thumb", Scale: 16},
			{Name: "thumb-gs", Scale: 16, Gs: true},
		}
	)

	//photos, err := takeouttoo.Find("/home/trimble/takeout01")
	//func FromFiles(jsonPath, resizePath string, sizes resize.Sizes) (photos entity.Photos, err error) {
	photos, err := takeouttoo.FromFiles(src)
	if err != nil {
		panic(err)
	}

	err = takeouttoo.AddResize(photos, dst, sizes)

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
