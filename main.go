package main

import (
	"bytes"
	"github.com/arjantop/imageoptimizer/ssim"
	"github.com/jszwec/csvutil"
	"github.com/nfnt/resize"
	"github.com/will-lol/walker"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

type Result struct {
	Resolution int
	Bytes      int
}

func main() {
	const inputDir = "./images/"
	const csvDir = "./csv/"
	resolutions := [...]int{2448, 2295, 2142, 1989, 1836, 1683, 1530, 1377, 1224, 1071, 918, 765, 612, 459, 306, 153}
	const targetSSIM = 0.9

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	for i, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".png" {
			result := make([]Result, 0, len(resolutions))

			file, err := os.Open(inputDir + entry.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			image, err := png.Decode(file)
			if err != nil {
				log.Fatal(err)
			}

			for _, entry := range resolutions {
				resized := resize.Resize(uint(entry), 0, image, resize.NearestNeighbor)
				bytes, err := encodeToQuality(resized, targetSSIM, 0.01)
				log.Println(bytes)
				if err != nil {
					log.Fatal(err)
				}
				result = append(result, Result{Resolution: entry, Bytes: len(bytes)})
			}
			csvString, err := csvutil.Marshal(result)
			if err != nil {
				log.Fatal(err)
			}
			csvBytes := []byte(csvString)
			err = os.WriteFile(csvDir+strconv.Itoa(i)+".csv", csvBytes, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func encodeToQuality(img image.Image, quality float64, tolerance float64) ([]byte, error) {
	var buf bytes.Buffer
	currentSsim := 0.0
	walkerInstance := walker.NewWalker(3)
	jpegQuality := 40

	for {
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality})
		log.Println(jpegQuality)
		if err != nil {
			return nil, err
		}
		imgEncoded, err := jpeg.Decode(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return nil, err
		}
		grayscaleImg := toGray(img)
		grayscaleImgEncoded := toGray(imgEncoded)
		currentSsim = ssim.Ssim(&grayscaleImg, &grayscaleImgEncoded)
		walkerInstance.Walk(jpegQuality)
		if walkerInstance.Check(isEqual) {
			break
		}
		if isEqualWithinTolerance(currentSsim, quality, tolerance) {
			break
		}
		if currentSsim > quality {
			jpegQuality--
		} else if currentSsim < quality {
			jpegQuality++
		}
		buf.Reset()
	}
	return buf.Bytes(), nil
}

func isEqual(x, y any) bool {
	if !(x == 0 || y == 0) {
		return false
	}
	return x == y
}

func isEqualWithinTolerance(x, y, tolerance float64) bool {
	return math.Abs(x-y) < tolerance
}

func toGray(img image.Image) image.Gray {
	newImage := image.NewGray(img.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			newImage.Set(x, y, img.At(x, y))
		}
	}
	return *newImage
}