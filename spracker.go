package spracker

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

// Generate an array of SCSS variable definitions.
func GenerateScssVariables(folder string, sheetName string, sheetImg image.Image, sprites []Sprite) string {
	variables := make([]string, 0, len(sprites) + 2)

	sheetUrl := fmt.Sprintf("$%s-url: image-url(\"%s.png\");\n", sheetName, filepath.Join(folder, sheetName))
	variables = append(variables, sheetUrl)

	sheetWidth  := fmt.Sprintf("$%s-width: %#vpx;", sheetName, sheetImg.Bounds().Max.X - sheetImg.Bounds().Min.X)
	sheetHeight := fmt.Sprintf("$%s-height: %#vpx;", sheetName, sheetImg.Bounds().Max.Y - sheetImg.Bounds().Min.Y)
	def         := fmt.Sprintf("%s\n%s\n", sheetWidth, sheetHeight)
	variables = append(variables, def)

	for _, s := range sprites {
		isMag, name, factor := IsMagnified(s.Name)
		if !isMag {
			factor = 1
		}
		prefix := fmt.Sprintf("%s-%s", sheetName, name)
		vX     := fmt.Sprintf("$%s-x: %#vpx;", prefix, float32(-s.Min.X)/factor)
		vY     := fmt.Sprintf("$%s-y: %#vpx;", prefix, float32(-s.Min.Y)/factor)
		vW     := fmt.Sprintf("$%s-width: %#vpx;", prefix, float32(s.Width())/factor)
		vH     := fmt.Sprintf("$%s-height: %#vpx;", prefix, float32(s.Height())/factor)
		def    := fmt.Sprintf("%s\n%s\n%s\n%s\n", vX, vY, vW, vH)
		variables = append(variables, def)
	}

	return strings.Join(variables, "\n")
}

// Generate an array of SCSS mixin definitions.
const mixinFormat string =
`@mixin %s-%s() {
  background: image-url("%s.png") no-repeat %#vpx %#vpx;%s
  width: %#vpx;
  height: %#vpx;
}
`
func GenerateScssMixins(folder string, sheetName string, sheetImg image.Image, sprites []Sprite) string {
	mixins := make([]string, 0, len(sprites))

	for _, s := range sprites {
		isMag, name, factor := IsMagnified(s.Name)
		bgSize := ""
		if !isMag {
			factor = 1
		} else {
			bgSize = fmt.Sprintf("\n  @include background-size(%vpx %vpx);", float32(sheetImg.Bounds().Max.X)/factor, float32(sheetImg.Bounds().Max.Y)/factor)
		}
		def := fmt.Sprintf(mixinFormat, sheetName, name, filepath.Join(folder, sheetName), float32(-s.Min.X)/factor, float32(-s.Min.Y)/factor, bgSize, float32(s.Width())/factor, float32(s.Height())/factor)
		mixins = append(mixins, def)
	}

	return strings.Join(mixins, "\n")
}

// Generate an array of CSS class definitions.
const classFormat string =
`.%s-%s {
  background: image-url("%s.png") no-repeat %#vpx %#vpx;%s
  width: %#vpx;
  height: %#vpx;
}
`
func GenerateCssClasses(folder string, sheetName string, sheetImg image.Image, sprites []Sprite) string {
	classes := make([]string, 0, len(sprites))

	for _, s := range sprites {
		isMag, name, factor := IsMagnified(s.Name)
		bgSize := ""
		if !isMag {
			factor = 1
		} else {
			bgSize = fmt.Sprintf("\n  @include background-size(%vpx %vpx);", float32(sheetImg.Bounds().Max.X)/factor, float32(sheetImg.Bounds().Max.Y)/factor)
		}
		class := fmt.Sprintf(classFormat, sheetName, name, filepath.Join(folder, sheetName), float32(-s.Min.X)/factor, float32(-s.Min.Y)/factor, bgSize, float32(s.Width())/factor, float32(s.Height())/factor)
		classes = append(classes, class)
	}

	return strings.Join(classes, "\n")
}

func SpritesModified(folder string) (bool, error) {
	sheetFileName := folder + ".png"
	sheetStat, err := os.Stat(sheetFileName)
	if err != nil {
		return true, nil // basically checking to see if the spritesheet exists
	}
	sheetModTime := sheetStat.ModTime()
	folderStat, err := os.Stat(folder)
	if err != nil {
		return false, err
	}
	folderModTime := folderStat.ModTime()
	if folderModTime.After(sheetModTime) {
		return true, nil
	}
	dir, err := os.Open(folder)
	if err != nil {
		return false, err
	}
	imageNames, err := dir.Readdirnames(0)
	if err != nil {
		return false, err
	}
	for _, imgBaseName := range imageNames {
		imgPath := filepath.Join(folder, imgBaseName)
		imgStat, err := os.Stat(imgPath)
		if err != nil {
			continue
		}
		imgModTime := imgStat.ModTime()
		if imgModTime.After(sheetModTime) {
			return true, nil
		}
	}
	return false, nil
}

// Read a folder containing subfolders which contain sprites, and generate a
// spritesheet for each subfolder.
func GenerateSpriteSheetsFromFolders(superFolder string, outputFolder string, checkTimeStamps bool, log *golog.Logger) (spriteSheets []Image, styleSheets []string, err error) {
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

	var anyErrors error = nil
	for _, folder := range folders {
		if checkTimeStamps {
			modified, err := SpritesModified(filepath.Join(superFolder, folder))
			if err != nil {
				anyErrors = err
				continue
			}
			if !modified {
				continue
			}
		}
		images, err := ReadImageFolder(filepath.Join(superFolder, folder), log)
		if err != nil {
			anyErrors = err
			continue
		}
		sheet, sprites := GenerateSpriteSheet(images)
		vars    := GenerateScssVariables(outputFolder, folder, sheet, sprites)
		mixins  := GenerateScssMixins(outputFolder, folder, sheet, sprites)
		classes := GenerateCssClasses(outputFolder, folder, sheet, sprites)
		spriteSheets = append(spriteSheets, Image{folder, sheet})
		styleSheets  = append(styleSheets, fmt.Sprintf("%s\n%s\n%s", vars, mixins, classes))
	}
	err = anyErrors

	return
}

// Check whether a sprite name describes a hi-res (e.g. retina) sprite. If so,
// return the base name and magnification factor.
func IsMagnified(name string) (isIt bool, baseName string, factor float32) {

	segments := strings.Split(name, "@")
	if len(segments) != 2 {
		return false, name, 1
	}
	baseName = segments[0]
	factorStr := segments[1]
	if factorStr[len(factorStr)-1] == 'x' {
		f, err := strconv.ParseFloat(factorStr[0:len(factorStr)-1], 32)
		if err != nil {
			return false, name, 1
		}
		factor = float32(f)
		isIt = true
	} else {
		return false, name, 1
	}
	return
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

// Write a stylesheet string to a file.
func WriteStyleSheet(style string, folder string, name string, log *golog.Logger) (err error) {
	fullname := filepath.Join(folder, name + ".scss")
	os.Remove(fullname)
	outFile, err := os.Create(fullname)
	if err != nil {
		log.Warning("Error opening SCSS output file %s", fullname)
		return err
	}
	_, err = outFile.WriteString(style)
	if err != nil {
		log.Warning("Error outputting SCSS to %s", fullname)
		return err
	}
	outFile.Close()
	return nil
}