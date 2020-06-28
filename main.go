package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var validExtensions = map[string]bool{}

func getImageNamesInDirectory(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, len(files))
	var ext string
	for idx, f := range files {
		ext = strings.ToLower(filepath.Ext(f.Name()))

		if _,ok := validExtensions[ext]; ok {
			names[idx] = f.Name()
		}
	}

	return names
}

func resizeImage(inPath, outPath, name string, maxWidth int, maxHeight int) error {
	f,err := os.Open(inPath + "/" + name)
	if err != nil {
		return err
	}
	defer f.Close()

	img,err := jpeg.Decode(f)
	if err != nil {
		return err
	}


	imgResized := resize.Thumbnail(uint(maxWidth), uint(maxHeight), img, resize.MitchellNetravali)

	outfile, err := os.Create(outPath + "/" + name)
	if err != nil {
		return err
	}
	defer outfile.Close()
	switch filepath.Ext(name) {
	case ".jpg":
		err = jpeg.Encode(outfile, imgResized, &jpeg.Options{
			Quality: 100,
		})
	case ".png":
		err = png.Encode(outfile, imgResized)
	default:
		return errors.New("unsupported file extension")
	}
	return err
}

func main() {
	validExtensions[".jpg"] = true
	validExtensions[".png"] = true

	inputPtr := flag.String("dir", "", "Path to images")
	outputPtr := flag.String("output", "", "Path to output (must exist)")
	maxWidthPtr := flag.Int("maxWidth", -1, "Max width in pixels")
	maxHeightPtr := flag.Int("maxWidth", -1, "Max height in pixels")
	flag.Parse()
	if *inputPtr == "" {
		panic("No input directory specified. Please run with the flag -dir=\"<dir>\"")
	}
	if *outputPtr == "" {
		panic("No output directory specified. Please run with the flag -output=\"<dir>\"")
	}
	if *maxWidthPtr == -1{
		panic("Please specify a max Width with -maxWidth=X")
	}
	if *maxHeightPtr == -1 {
		panic("Please specify a max Height with -maxHeight=X")
	}

	images := getImageNamesInDirectory(*inputPtr)

	log.Println(fmt.Sprintf("%d images to process", len(images)))

	failed := make([]error, 0)
	failedMutex := sync.Mutex{}

	var err error
	wg := sync.WaitGroup{}
	wg.Add(len(images))
	for idx := range images {
		go func(path string) {
			err = resizeImage(*inputPtr, *outputPtr, path, *maxWidthPtr, *maxHeightPtr)
			if err != nil {
				failedMutex.Lock()
				failed = append(failed, err)
				failedMutex.Unlock()
			}
			log.Println(fmt.Sprintf("Processed: %s", path))
			wg.Done()
		}(images[idx])
	}

	wg.Wait()


	if len(failed) > 0 {
		log.Println("Finished, with errors:")
		for _,e := range failed {
			log.Println(e)
		}
	} else {
		log.Println("Finished")
	}
}
