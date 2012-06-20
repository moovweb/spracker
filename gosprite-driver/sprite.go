package main

import (
  "os"
  "fmt"
  "gosprite"
  "golog"
)

func main() {
  if len(os.Args) < 3 {
    fmt.Println("Usage: sprite [image folder] [output filename]")
    os.Exit(0)
  }
  path    := os.Args[1]
  outName := os.Args[2]
  log     := golog.NewLogger("sprite")

  images, _ := gosprite.ReadImageFolder(path, log)
  sheet, _  := gosprite.GenerateSpriteSheet(images, log)

  fmt.Printf("Read %d image files.\n", len(images))
  gosprite.WriteSpriteSheet(sheet, path, outName, log)
}