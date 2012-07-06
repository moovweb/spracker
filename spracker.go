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

// We need to bundle an image with it's name (derived from its filename).
type Image struct {
	Name string
	image.Image
}

// Bundle the name of a sprite image with its viewport into the spritesheet.
type Sprite struct {
	Name string
	image.Rectangle
}

// Helper methods.
func (s *Sprite) Width() int {
	return s.Max.X - s.Min.X
}

func (s *Sprite) Height() int {
	return s.Max.Y - s.Min.Y
}

// Take a folder name, read all the images within, and return them in an array.
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

// Compose the spritesheet image and return an array of individual sprite data.
func GenerateSpriteSheet(images []Image) (sheet draw.Image, sprites []Sprite) {
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

// Write the spritesheet image to a file.
func WriteSpriteSheet(img image.Image, folder string, name string, log *golog.Logger) (err error) {
	fullname := filepath.Join(folder, name + ".png")
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


// Generate an array of SCSS variable definitions.
func GenerateScssVariables(sheetName string, sheetImg image.Image, sprites []Sprite) string {
	variables := make([]string, 0, len(sprites) + 2)

	sheetUrl := fmt.Sprintf("$spritesheet-%s-url: image-url(\"%s.png\");\n", sheetName, sheetName)
	variables = append(variables, sheetUrl)

	sheetWidth  := fmt.Sprintf("$spritesheet-%s-width: %d;", sheetName, sheetImg.Bounds().Max.X - sheetImg.Bounds().Min.X)
	sheetHeight := fmt.Sprintf("$spritesheet-%s-height: %d;", sheetName, sheetImg.Bounds().Max.Y - sheetImg.Bounds().Min.Y)
	def         := fmt.Sprintf("%s\n%s\n", sheetWidth, sheetHeight)
	variables = append(variables, def)

	for _, s := range sprites {
		prefix := fmt.Sprintf("sprite-%s-%s", sheetName, s.Name)
		vX     := fmt.Sprintf("$%s-x: %dpx;", prefix, -s.Min.X)
		vY     := fmt.Sprintf("$%s-y: %dpx;", prefix, -s.Min.Y)
		vW     := fmt.Sprintf("$%s-width: %dpx;", prefix, s.Width())
		vH     := fmt.Sprintf("$%s-height: %dpx;", prefix, s.Height())
		def    := fmt.Sprintf("%s\n%s\n%s\n%s\n", vX, vY, vW, vH)
		variables = append(variables, def)
	}

	return strings.Join(variables, "\n")
}

// Generate an array of SCSS mixin definitions.
const mixinFormat string =
`@mixin sprite-%s-%s() {
  background: image-url("%s.png") no-repeat %dpx %dpx;
  width: %dpx;
  height: %dpx;
}
`
func GenerateScssMixins(sheetName string, sprites []Sprite) string {
	mixins := make([]string, 0, len(sprites))

	for _, s := range sprites {
		def := fmt.Sprintf(mixinFormat, sheetName, s.Name, sheetName, -s.Min.X, -s.Min.Y, s.Width(), s.Height())
		mixins = append(mixins, def)
	}

	return strings.Join(mixins, "\n")
}

// Read 
func GenerateSpriteSheetsFromFolders(superFolder string, log *golog.Logger) (spriteSheets []Image, styleSheets []string, err error) {
	container, err := os.Open(superFolder)
	if err != nil {
		log.Warning("Error opening folder %s containing sprite subfolders", superFolder)
		return nil, nil, err
	}
	folders, err := container.Readdirnames(0)
	container.Close()
	if err != nil {
		log.Warning("Error gathering names of sprite subfolders in %s", superFolder)
		return nil, nil, err
	}

	spriteSheets = make([]Image, 0)
	styleSheets  = make([]string, 0)

	for _, folder := range folders {
		images, err := ReadImageFolder(folder, log)
		if err != nil {
			return nil, nil, err // TO DO: make more resilient
		}
		sheet, sprites := GenerateSpriteSheet(images)
		vars   := GenerateScssVariables(folder, sheet, sprites)
		mixins := GenerateScssMixins(folder, sprites)
		spriteSheets = append(spriteSheets, Image{folder, sheet})
		styleSheets  = append(styleSheets, fmt.Sprintf("%s\n%s", vars, mixins))
	}

	return
}