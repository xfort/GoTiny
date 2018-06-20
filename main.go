package main

import (
	"github.com/xfort/GoTiny/tiny"
	"log"
)

func main() {

	tinyHandler := &tiny.TinyHandler{}
	tinyHandler.SetData("xx", "/home/tongying/doc/tmp/res/drawable-hdpi")
	err := tinyHandler.CompressAllImages("/home/tongying/doc/tmp/res_image/drawable-hdpi", "")
	if err != nil {
		log.Println(err)
	}
}
