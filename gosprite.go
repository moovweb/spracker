package gosprite

import (
	"os"
	"path/filepath"
	"image"
	"image/png"
	"golog"
)

func ReadImageFolder(path string, log *golog.Logger) (images []image.Image, err error) {
	dir, err := os.Open(path)
	if err != nil {
		log.Warning("Error opening sprite folder %s", path)
		return nil, err
	}
	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Warning("Error gathering names of sprite files in %s", path)
		return nil, err
	}

	images = make([]image.Image, 0, len(names))
	for _, name := range names {
		imgFile, err := os.Open(filepath.Join(path, name))
		defer imgFile.Close()
		if err != nil {
			log.Warning("Error reading sprite image from %s", name)
			continue
		}
		img, err := png.Decode(imgFile)
		if err != nil {
			log.Warning("Error decoding png sprite image in %s", name)
			continue
		}
		images = append(images, img)
	}

	return images, nil
}