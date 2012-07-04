package main

import (
	"os"
)

import (
	"golog"
	"spracker"
)

func main() {
	if len(os.Args) < 3 {
		println("Usage: sprite [image folder] [output filename]")
		os.Exit(0)
	}
	path := os.Args[1]
	outName := os.Args[2]
	log := golog.NewLogger("sprite")

	images, _ := spracker.ReadImageFolder(path, log)
	sheet, _ := spracker.GenerateSpriteSheet(images, log)

	println("Read", len(images), "image files:")
	for _, img := range images {
		println(img.Name)
	}
	spracker.WriteSpriteSheet(sheet, "", outName, log)
}
