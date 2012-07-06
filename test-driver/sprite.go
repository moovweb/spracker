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
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	log := golog.NewLogger("sprite")

	// images, _ := spracker.ReadImageFolder(path, log)
	// sheet, sprites := spracker.GenerateSpriteSheet(images)

	// // println("Read", len(images), "image files. Sprite layout:")
	// // for i, img := range images {
	// // 	s := sprites[i]
	// // 	fmt.Printf("%s: (%d,%d)\t-\t(%d,%d)\n", img.Name, s.Min.X, s.Min.Y, s.Max.X, s.Max.Y)
	// // }
	// spracker.WriteSpriteSheet(sheet, "", path, log)

	// println(spracker.GenerateScssVariables(path, sheet, sprites))
	// println(spracker.GenerateScssMixins(path, sprites))

	sheets, styles, _ := spracker.GenerateSpriteSheetsFromFolders(path, log)
	for i, sheet := range sheets {
		spracker.WriteSpriteSheet(sheet.Image, path, sheet.Name, log)
		spracker.WriteStyleSheet(styles[i], path, sheet.Name, log)
	}
}
