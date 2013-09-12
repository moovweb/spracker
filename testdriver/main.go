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
	os.RemoveAll("kittens")
	os.Mkdir("kittens", 0755)

	log := golog.NewLogger("")
	log.AddProcessor("first", golog.NewConsoleProcessor(golog.LOG_INFO, true))

	// Fetch ten pseudo-randomly sized images from http://placekitten.com and
	// write them to disk.
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())

		size := 64 + rand.Int()%20
		magFrac := float64(rand.Int()%2*5) / 10
		magFac := float64(rand.Int()%4+1) + magFrac
		size = int(math.Floor(magFac)) * (64 + rand.Int()%20)

		imgResp, netErr := http.Get(fmt.Sprintf("http://placekitten.com/%d/%d", size, size))
		if netErr != nil || imgResp.Status != "200 OK" {
			log.Errorf("unable to download test sprite of size %d x %d", size, size)
			continue
		}

		imgData, readErr := ioutil.ReadAll(imgResp.Body)
		if readErr != nil {
			log.Errorf("unable to read image data")
			continue
		}

		var fileName string
		if magFac != 1 {
			fileName = fmt.Sprintf("kittens/k0%d@%#vx.jpg", i, magFac)
		} else {
			fileName = fmt.Sprintf("kittens/k0%d.jpg", i)
		}

		writeErr := ioutil.WriteFile(fileName, imgData, 0666)
		if writeErr != nil {
			log.Errorf("unable to write image data")
		}
	}

	// wd, _ := os.Getwd()
	sheet, styles, _, spriteErr := spracker.GenerateSpriteSheetFromFolder("kittens", "kittens", ".", false, false, log)
	if spriteErr != nil {
		log.Errorf("unable to generate spritesheet; please delete images and stylesheets and try again")
	} else {
		wspError := spracker.WriteSpriteSheet(sheet.Image, "kittens", sheet.Name, log)
		wstError := spracker.WriteStyleSheet(styles, "kittens", sheet.Name+".css", log)
		if wspError != nil {
			log.Errorf("unable to write spritesheet")
		}
		if wstError != nil {
			log.Errorf("unable to write stylesheet")
		}
	}

	golog.FlushLogsAndDie()
}
