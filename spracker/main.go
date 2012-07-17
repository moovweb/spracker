package main

import (
	"flag"
	// "fmt"
	// "os"
)

import (
	"golog"
	"spracker"
)

func main() {

	var (
		generateScss       bool
		spritesFolder      string
		spriteSheetsFolder string
		styleSheetsFolder  string
	)

	flag.BoolVar(&generateScss, "scss", true, "generate Sass/SCSS variables and mixins")
	flag.StringVar(&spritesFolder, "sprites-folder", ".", "folder containing subfolders with sprite images")
	flag.StringVar(&spriteSheetsFolder, "spritesheets-folder", ".", "output folder in which to deposit the sprite-sheets")
	flag.StringVar(&styleSheetsFolder, "stylesheets-folder", ".", "output folder in which to deposit the stylesheets")

	flag.Parse()

	log := golog.NewLogger("spracker")

	sheets, styles, _ := spracker.GenerateSpriteSheetsFromFolders(spritesFolder, spriteSheetsFolder, true, log)
	for i, sheet := range sheets {
		spracker.WriteSpriteSheet(sheet.Image, spriteSheetsFolder, sheet.Name, log)
		spracker.WriteStyleSheet(styles[i], styleSheetsFolder, sheet.Name, log)
	}

}
