package gosprite

import (
	"os"
	"path/filepath"
	"image"
	"image/png"
	"image/draw"
	"golog"
)

func ReadImageFolder(path string, log *golog.Logger) (images map[string]image.Image, err error) {
	dir, err := os.Open(path)
	if err != nil {
		log.Warning("Error opening sprite folder %s", path)
		return nil, err
	}
	names, err := dir.Readdirnames(0)
	dir.Close()
	if err != nil {
		log.Warning("Error gathering names of sprite files in %s", path)
		return nil, err
	}

	images = make(map[string]image.Image)
	for _, name := range names {
		imgFile, err := os.Open(filepath.Join(path, name))
		if err != nil {
			log.Warning("Error reading sprite image from %s", name)
			continue
		}
		img, err := png.Decode(imgFile)
		if err != nil {
			log.Warning("Error decoding png sprite image in %s", name)
			continue
		}
		imgFile.Close()
		images[name] = img
	}

	return images, nil
}

func GenerateSpriteSheet(images map[string]image.Image, log *golog.Logger) (sheet draw.Image, sprites map[string]image.Rectangle) {
	var sheetHeight int = 0
	var sheetWidth  int = 0
	sprites = make(map[string]image.Rectangle)

	// calculate the size of the spritesheet and accumulate data for the
	// individual sprites
	for name, img := range images {
		bounds := img.Bounds()
		sprites[name] = image.Rect(0, sheetHeight, bounds.Dx(), sheetHeight + bounds.Dy())
		sheetHeight += bounds.Dy()
		if bounds.Dx() > sheetWidth {
			sheetWidth = bounds.Dx()
		}
	}

	// create the sheet image
	sheet = image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))

	// compose the sheet
	for name, img := range images {
		draw.Draw(sheet, sprites[name], img, image.Pt(0, 0), draw.Src)
	}

	return sheet, sprites
}

func WriteSpriteSheet(img image.Image, path string, name string, log *golog.Logger) (err error) {
	fullname := filepath.Join(path, name)
	os.Remove(fullname)
	outFile, err := os.Create(fullname)
	if err != nil {
		log.Warning("Error opening spritesheet output file %s", fullname)
		return err
	}
	err = png.Encode(outFile, img)
	if err != nil {
		log.Warning("Error outputting spritesheet to %s", fullname)
		return err
	}
	outFile.Close()
	return nil
}