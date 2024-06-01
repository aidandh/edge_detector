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
	filter :=
		[][]int{
			{-1, 0, -1},
			{0, 4, 0},
			{-1, 0, -1},
		}
	filterHeight := len(filter)
	filterWidth := len(filter[0])
	imageWidth := original.Bounds().Max.X
	imageHeight := original.Bounds().Max.Y
	laplacian := image.NewRGBA64(original.Bounds())

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			lr, lg, lb := 0, 0, 0
			_, _, _, la := original.At(x, y).RGBA()
			for iHeight := 0; iHeight < filterHeight; iHeight++ {
				for iWidth := 0; iWidth < filterWidth; iWidth++ {
					xCoord := (iWidth - filterWidth/2 + iWidth + imageWidth) % imageWidth
					yCoord := (iHeight - filterHeight/2 + iHeight + imageHeight) % imageHeight
					// fmt.Println(iWidth, iHeight)
					or, og, ob, _ := original.At(xCoord, yCoord).RGBA()
					lr += int(or) * filter[iHeight][iWidth]
					lg += int(og) * filter[iHeight][iWidth]
					lb += int(ob) * filter[iHeight][iWidth]
					// fmt.Println("lr:", lr, "or:", or, "filter[iHeight][iWidth]:", filter[iHeight][iWidth])
					// fmt.Println("lg:", lg, "og:", og, "filter[iHeight][iWidth]:", filter[iHeight][iWidth])
					// fmt.Println("lb:", lb, "ob:", ob, "filter[iHeight][iWidth]:", filter[iHeight][iWidth])
				}
			}
			laplacian.Set(x, y, color.RGBA64{
				uint16(lr),
				uint16(lg),
				uint16(lb),
				uint16(la),
			})
			// r, g, b, a := original.At(x, y).RGBA()
			// fmt.Println(x, y)
			// fmt.Println(r, g, b, a)
			// fmt.Println(uint16(lr), uint16(lg), uint16(lb), uint16(la))
			// fmt.Println("---------------------------")
			// time.Sleep(500)
			time.Sleep(0)
		}
	}

	return laplacian
}
