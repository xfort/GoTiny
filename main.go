package main

import (
	"github.com/xfort/GoTiny/tiny"
	"log"
	"path/filepath"
	"os"
)

func main() {
	//mask := syscall.Umask(0)
	//defer syscall.Umask(mask)

	imgsSir := "/home/tongying/doc/TongYe/doc/icons/tmp/drawable-xxhdpi"
	outDir := filepath.Join(filepath.Dir(imgsSir), "out", filepath.Base(imgsSir))
	err := os.MkdirAll(outDir, 0766)
	if err != nil {
		log.Fatalln(err)
	}
	tinyHandler := &tiny.TinyHandler{}
	tinyHandler.SetData("vX2I6B4hbbwwffnZwx78qS4OsrG6JB4L", outDir)
	err = tinyHandler.CompressAllImages(imgsSir, "")
	if err != nil {
		log.Println(err)
	}
}
