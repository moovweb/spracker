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
	if len(os.Args) < 2 {
		println("Usage: sprite [image folder]")
		os.Exit(0)
	}
	path := os.Args[1]
	log := golog.NewLogger("sprite")

	images, _ := spracker.ReadImageFolder(path, log)
	sheet, sprites := spracker.GenerateSpriteSheet(images)

	// println("Read", len(images), "image files. Sprite layout:")
	// for i, img := range images {
	// 	s := sprites[i]
	// 	fmt.Printf("%s: (%d,%d)\t-\t(%d,%d)\n", img.Name, s.Min.X, s.Min.Y, s.Max.X, s.Max.Y)
	// }
	spracker.WriteSpriteSheet(sheet, "", path, log)

	println(spracker.GenerateScssVariables(path, sheet, sprites))
	println(spracker.GenerateScssMixins(path, sprites))
}
