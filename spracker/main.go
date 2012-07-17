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
		checkTimestamps    bool
		spritesFolder      string
		spriteSheetsFolder string
		styleSheetsFolder  string
	)

	flag.BoolVar(&generateScss, "scss", true, "generate Sass/SCSS variables and mixins")
	flag.BoolVar(&checkTimestamps, "check-timestamps", true, "don't regenerate sprite-sheets if they're newer than their component sprite images")
	flag.StringVar(&spritesFolder, "sprites-folder", ".", "input folder containing subfolders with sprite images")
	flag.StringVar(&spriteSheetsFolder, "spritesheets-folder", ".", "output folder in which to deposit the sprite-sheets")
	flag.StringVar(&styleSheetsFolder, "stylesheets-folder", ".", "output folder in which to deposit the stylesheets")

	flag.Parse()

	log := golog.NewLogger("")
	log.AddProcessor("first", golog.NewConsoleProcessor(golog.LOG_INFO, true))

	var stylesheetExtension string;
	if (generateScss) {
		stylesheetExtension = ".scss"
	} else {
		stylesheetExtension = ".css"
	}

	sheets, styles, gErr := spracker.GenerateSpriteSheetsFromFolders(spritesFolder, spriteSheetsFolder, generateScss, checkTimestamps, log)
	var wErr error
	for i, sheet := range sheets {
		wspErr := spracker.WriteSpriteSheet(sheet.Image, spriteSheetsFolder, sheet.Name, log)
		wstErr := spracker.WriteStyleSheet(styles[i], styleSheetsFolder, sheet.Name + stylesheetExtension, log)
		if wspErr != nil || wstErr != nil {
			wErr = wstErr
		}
	}

	if gErr != nil || wErr != nil {
		golog.FlushLogsAndDie()
	}

}
