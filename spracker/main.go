package main

import (
	"flag"
	"path/filepath"
	// "fmt"
	// "os"
)

import (
	"github.com/moovweb/golog"
	"github.com/moovweb/spracker"
)

func main() {

	var (
		generateScss    bool
		checkTimestamps bool
		// projectFolder      string
		spritesFolder      string
		spriteSheetsFolder string
		styleSheetsFolder  string
	)

	flag.BoolVar(&generateScss, "scss", true, "generate Sass/SCSS variables and mixins")
	flag.BoolVar(&checkTimestamps, "check-timestamps", true, "don't regenerate sprite-sheets if they're newer than their component sprite images")
	// flag.StringVar(&projectFolder, "project-folder", ".", "base folder for your project")
	flag.StringVar(&spritesFolder, "sprites-folder", ".", "input folder containing subfolders with sprite images")
	flag.StringVar(&spriteSheetsFolder, "spritesheets-folder", ".", "output folder in which to deposit the sprite-sheets")
	flag.StringVar(&styleSheetsFolder, "stylesheets-folder", ".", "output folder in which to deposit the stylesheets")

	flag.Parse()

	log := golog.NewLogger("")
	log.AddProcessor("first", golog.NewConsoleProcessor(golog.LOG_INFO, true))

	var stylesheetExtension string
	if generateScss {
		stylesheetExtension = ".scss"
	} else {
		stylesheetExtension = ".css"
	}

	sheets, styles, _ := spracker.GenerateSpriteSheetsFromFolders(spritesFolder, spriteSheetsFolder, spriteSheetsFolder, generateScss, checkTimestamps, log)
	for i, sheet := range sheets {
		wstErr := spracker.WriteSpriteSheet(sheet.Image, spriteSheetsFolder, sheet.Name, log)
		if wstErr == nil {
			log.Infof("Generated sprite-sheet '%s.png'", filepath.Join(spriteSheetsFolder, sheet.Name))
		}
		wspErr := spracker.WriteStyleSheet(styles[i], styleSheetsFolder, sheet.Name+stylesheetExtension, log)
		if wspErr == nil {
			log.Infof("Generated stylesheet '%s%s'", filepath.Join(styleSheetsFolder, sheet.Name), stylesheetExtension)
		}
	}

	golog.FlushLogsAndDie()
}
