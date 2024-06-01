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
)

type ImageWithName struct {
	Image image.Image
	Name  string
}

func main() {
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
	for _, image := range images {
		laplacianImage := applyLaplacianFilter(image.Image)

		outputFile, err := os.Create("output/" + image.Name + ".png")
		if err != nil && !os.IsExist(err) {
			fmt.Println(err.Error())
			continue
		}
		defer outputFile.Close()
		png.Encode(outputFile, laplacianImage)
	}
}

func openImages(paths []string) []ImageWithName {
	images := make([]ImageWithName, 0, len(paths))
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer file.Close()

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
	}
	return images
}

func applyLaplacianFilter(original image.Image) image.Image {
	laplacian := image.NewRGBA(original.Bounds())

	for x := range laplacian.Bounds().Max.X - 1 {
		for y := range laplacian.Bounds().Max.Y - 1 {
			r, g, b, a := original.At(x, y).RGBA()
			laplacian.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return laplacian
}
