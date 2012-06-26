package main

import (
	"fmt"
	"os"
)

import (
	"golog"
	"spracker"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sprite [image folder] [output filename]")
		os.Exit(0)
	}
	path := os.Args[1]
	outName := os.Args[2]
	log := golog.NewLogger("sprite")

	images, _ := spracker.ReadImageFolder(path, log)
	sheet, _ := spracker.GenerateSpriteSheet(images, log)

	fmt.Printf("Read %d image files.\n", len(images))
	spracker.WriteSpriteSheet(sheet, path, outName, log)
}
