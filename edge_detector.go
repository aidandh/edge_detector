package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"
)

func main() {
	err := os.Mkdir("output", 0750)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err.Error())
		return
	}

	paths := os.Args[1:]
	files := openFiles(paths)
	for _, file := range files {
		defer file.Close()
		image, err := applyFilter(file)
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
		outputFileName := strings.Split(splitStr[len(splitStr)-1], ".")[0]
		outputFile, err := os.Create("output/" + outputFileName + ".png")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer outputFile.Close()
		png.Encode(outputFile, image)
	}
}

func openFiles(paths []string) []*os.File {
	files := make([]*os.File, 0, len(paths))
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		files = append(files, file)
	}
	return files
}

func applyFilter(file *os.File) (image.Image, error) {
	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return image, nil
}
