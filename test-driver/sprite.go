package main

import (
	// "fmt"
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
	sheet, sprites := spracker.GenerateSpriteSheet(images, log)

	// println("Read", len(images), "image files. Sprite layout:")
	// for i, img := range images {
	// 	s := sprites[i]
	// 	fmt.Printf("%s: (%d,%d)\t-\t(%d,%d)\n", img.Name, s.Min.X, s.Min.Y, s.Max.X, s.Max.Y)
	// }
	spracker.WriteSpriteSheet(sheet, "", outName, log)

	variables := spracker.GenerateScssVariables(path, sheet, sprites)
	for _, v := range variables {
		println(v)
	}

	mixins := spracker.GenerateScssMixins(path, sprites)
	for _, m := range mixins {
		println(m)
	}
}
