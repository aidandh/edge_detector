package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"
)

type ImageWithName struct {
	Data image.Image
	Name string
}

func main() {
	startTime := time.Now()
	err := os.Mkdir("output", 0750)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err.Error())
		return
	}

	paths := os.Args[1:]
	if len(paths) == 0 {
		fmt.Println("Too few arguments")
		return
	}

	images := openImages(paths)
	chans := make([]chan image.Image, 0, len(images))

	for _, curImage := range images {
		c := make(chan image.Image)
		chans = append(chans, c)
		go applyLaplacianFilter(curImage.Data, c)
	}

	for i, c := range chans {
		outputFile, err := os.Create("output/" + images[i].Name + ".png")
		if err != nil && !os.IsExist(err) {
			fmt.Println(err.Error())
			continue
		}
		defer outputFile.Close()
		laplacianImage := <-c
		png.Encode(outputFile, laplacianImage)
	}

	duration := time.Since(startTime)
	fmt.Println("Execution time:", duration)
}

func openImages(paths []string) []ImageWithName {
	images := make([]ImageWithName, 0, len(paths))
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		image, _, err := image.Decode(file)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		var splitStr []string
		if strings.Contains(file.Name(), "/") {
			splitStr = strings.Split(file.Name(), "/")
		} else {
			splitStr = strings.Split(file.Name(), "\\")
		}
		imageName := strings.Split(splitStr[len(splitStr)-1], ".")[0]

		images = append(images, ImageWithName{image, imageName})

		file.Close()
	}
	return images
}

func applyLaplacianFilter(original image.Image, c chan image.Image) {
	defer close(c)
	filter :=
		[][]int{
			{-1, -1, -1},
			{-1, 8, -1},
			{-1, -1, -1},
		}
	filterHeight := len(filter)
	filterWidth := len(filter[0])
	imageWidth := original.Bounds().Max.X
	imageHeight := original.Bounds().Max.Y
	laplacian := image.NewRGBA(original.Bounds())

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			lr, lg, lb := 0, 0, 0
			_, _, _, la := original.At(x, y).RGBA()
			for iHeight := 0; iHeight < filterHeight; iHeight++ {
				for iWidth := 0; iWidth < filterWidth; iWidth++ {
					xCoord := (x - filterWidth/2 + iWidth + imageWidth) % imageWidth
					yCoord := (y - filterHeight/2 + iHeight + imageHeight) % imageHeight
					or, og, ob, _ := original.At(xCoord, yCoord).RGBA()
					lr += int(or/256) * filter[iHeight][iWidth]
					lg += int(og/256) * filter[iHeight][iWidth]
					lb += int(ob/256) * filter[iHeight][iWidth]
				}
			}
			laplacian.Set(x, y, color.RGBA{
				uint8(clamp(lr)),
				uint8(clamp(lg)),
				uint8(clamp(lb)),
				uint8(la),
			})
		}
	}

	c <- laplacian
}

func clamp(value int) int {
	if value < 0 {
		return 0
	} else if value > 255 {
		return 255
	} else {
		return value
	}
}
