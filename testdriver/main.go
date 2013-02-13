package main

import (
  "fmt"
  "io/ioutil"
  "math"
  "math/rand"
  "net/http"
  "os"
  "time"
)

import (
  "golog"
  "spracker"
)

func main() {

  // var (
  //   generateScss    bool
  //   checkTimestamps bool
  //   spritesFolder      string
  //   spriteSheetsFolder string
  //   styleSheetsFolder  string
  // )

  // flag.BoolVar(&generateScss, "scss", true, "generate Sass/SCSS variables and mixins")
  // flag.BoolVar(&checkTimestamps, "check-timestamps", true, "don't regenerate sprite-sheets if they're newer than their component sprite images")
  // flag.StringVar(&spritesFolder, "sprites-folder", ".", "input folder containing subfolders with sprite images")
  // flag.StringVar(&spriteSheetsFolder, "spritesheets-folder", ".", "output folder in which to deposit the sprite-sheets")
  // flag.StringVar(&styleSheetsFolder, "stylesheets-folder", ".", "output folder in which to deposit the stylesheets")

  // flag.Parse()


  // var stylesheetExtension string
  // if generateScss {
  //   stylesheetExtension = ".scss"
  // } else {
  //   stylesheetExtension = ".css"
  // }

  // sheets, styles, _ := spracker.GenerateSpriteSheetsFromFolders(spritesFolder, spriteSheetsFolder, spriteSheetsFolder, generateScss, checkTimestamps, log)
  // for i, sheet := range sheets {
  //   wstErr := spracker.WriteSpriteSheet(sheet.Image, spriteSheetsFolder, sheet.Name, log)
  //   if wstErr == nil {
  //     log.Info("Generated sprite-sheet '%s.png'", filepath.Join(spriteSheetsFolder, sheet.Name))
  //   }
  //   wspErr := spracker.WriteStyleSheet(styles[i], styleSheetsFolder, sheet.Name+stylesheetExtension, log)
  //   if wspErr == nil {
  //     log.Info("Generated stylesheet '%s%s'", filepath.Join(styleSheetsFolder, sheet.Name), stylesheetExtension)
  //   }
  // }


  log := golog.NewLogger("")
  log.AddProcessor("first", golog.NewConsoleProcessor(golog.LOG_INFO, true))

  // Fetch ten pseudo-randomly sized images from http://placekitten.com and
  // write them to disk.
  for i := 0; i < 10; i++ {
    rand.Seed(time.Now().UnixNano())

    size    := 64 + rand.Int() % 20
    magFrac := float64(rand.Int() % 2 * 5) / 10
    magFac  := float64(rand.Int() % 4 + 1) + magFrac
    size     = int(math.Floor(magFac)) * (64 + rand.Int() % 20)

    imgResp, netErr := http.Get(fmt.Sprintf("http://placekitten.com/%d/%d", size, size))
    if netErr != nil || imgResp.Status != "200 OK" {
      log.Error("unable to download test sprite of size %d x %d", size, size)
      continue
    }

    imgData, readErr := ioutil.ReadAll(imgResp.Body)
    if readErr != nil {
      log.Error("unable to read image data")
      continue
    }

    var fileName string
    if magFac != 1 {
      fileName = fmt.Sprintf("k0%d@%#vx.jpg", i, magFac)
    } else {
      fileName = fmt.Sprintf("k0%d.jpg", i)
    }

    writeErr := ioutil.WriteFile(fileName, imgData, 0666)
    if writeErr != nil {
      log.Error("unable to write image data")
    }
  }

  wd, _ := os.Getwd()
  sheet, styles, _, spriteErr := spracker.GenerateSpriteSheetFromFolder(wd, "kittens", ".", false, false, log)
  if spriteErr != nil {
    log.Error("unable to generate spritesheet; please delete images and stylesheets and try again")
  } else {
    wspError := spracker.WriteSpriteSheet(sheet.Image, "kittens", sheet.Name, log)
    wstError := spracker.WriteStyleSheet(styles, "kittens", sheet.Name+".css", log)
    if wspError != nil {
      log.Error("unable to write spritesheet")
    }
    if wstError != nil {
      log.Error("unable to write stylesheet")
    }
  }

  golog.FlushLogsAndDie()
}
