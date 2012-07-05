package spracker

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
)
import (
	"golog"
)

type Image struct {
	Name string
	image.Image
}

type Sprite struct {
	Name string
	image.Rectangle
}

func ReadImageFolder(path string, log *golog.Logger) (images []Image, err error) {
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

	// put these in a canonical order, otherwise it'll be platform-specific
	sort.Strings(names)

	images = make([]Image, 0)
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
		if strings.HasSuffix(name, ".png") {
			name = name[0 : len(name)-4]
		}
		images = append(images, Image{name, img})
	}

	return images, nil
}

func GenerateSpriteSheet(images []Image, log *golog.Logger) (sheet draw.Image, sprites []Sprite) {
	var (
		sheetHeight int = 0
		sheetWidth  int = 0
	)
	sprites = make([]Sprite, 0)

	// calculate the size of the spritesheet and accumulate data for the
	// individual sprites
	for _, img := range images {
		bounds := img.Bounds()
		sprites = append(sprites, Sprite{img.Name, image.Rect(0, sheetHeight, bounds.Dx(), sheetHeight+bounds.Dy())})
		sheetHeight += bounds.Dy()
		if bounds.Dx() > sheetWidth {
			sheetWidth = bounds.Dx()
		}
	}

	// create the sheet image
	sheet = image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))

	// compose the sheet
	for i, img := range images {
		draw.Draw(sheet, sprites[i].Rectangle, img, image.Pt(0, 0), draw.Src)
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

func GenerateScssVariables(sheetName string, sheetImg image.Image, sprites []Sprite) (variables []string) {
	variables = make([]string, 0, len(sprites) + 1)

	sheetWidth  := fmt.Sprintf("$spritesheet-%s-width: %d;", sheetName, sheetImg.Bounds().Max.X - sheetImg.Bounds().Min.X)
	sheetHeight := fmt.Sprintf("$spritesheet-%s-height: %d;", sheetName, sheetImg.Bounds().Max.Y - sheetImg.Bounds().Min.Y)
	def         := fmt.Sprintf("%s\n%s\n", sheetWidth, sheetHeight)
	variables = append(variables, def)

	for _, s := range sprites {
		prefix := fmt.Sprintf("sprite-%s-%s", sheetName, s.Name)
		vX     := fmt.Sprintf("$%s-x: %dpx;", prefix, s.Min.X)
		vY     := fmt.Sprintf("$%s-y: %dpx;", prefix, s.Min.Y)
		vW     := fmt.Sprintf("$%s-width: %dpx;", prefix, s.Max.X - s.Min.X)
		vH     := fmt.Sprintf("$%s-height: %dpx;", prefix, s.Max.Y - s.Min.Y)
		def    := fmt.Sprintf("%s\n%s\n%s\n%s\n", vX, vY, vW, vH)
		variables = append(variables, def)
	}

	return
}