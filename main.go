package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/llgcode/draw2d/draw2dimg"
	dotimage "github.com/tejesh-kaliki/dot-image-generator/src"
)

const (
	inputDir  = "test/input"
	outputDir = "test/output3"
)

func processImage(fileName, inputDir, outputDir string) error {
	imageFile, err := os.Open(inputDir + string(os.PathSeparator) + fileName)
	if err != nil {
		log.Printf("Error opening file %s: %v\n", fileName, err)
		return err
	}
	defer imageFile.Close()

	imageData, _, err := image.Decode(imageFile)
	if err != nil {
		log.Printf("Error reading image content %s: %v\n", fileName, err)
		return err
	}

	blockSize := uint(12)
	dotImageColors := dotimage.ComputeDotImageColors(imageData, blockSize)

	destImage := dotimage.CreateDotStyleImage(dotImageColors, uint(24))

	fileNameWithoutTag, _ := strings.CutSuffix(fileName, ".jpg")
	fileNameWithPng := fileNameWithoutTag + ".png"
	err = draw2dimg.SaveToPngFile(outputDir+string(os.PathSeparator)+fileNameWithPng, destImage)
	if err != nil {
		log.Printf("Error saving to file %s: %v\n", fileName, err)
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Please provide 2 parameters: input folder and output folder")
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	if inputDir == outputDir {
		fmt.Println("The input and output directories are set to same. This can override the input files.")
		fmt.Print("Are you sure? [yes|y|no|n] ")
		var userInput string
		fmt.Scanf("%s", &userInput)
		switch strings.ToLower(userInput) {
		case "no":
		case "n":
			os.Exit(0)
		case "yes":
		case "y":
			break
		default:
			fmt.Println("Invaild choice")
			os.Exit(1)
		}
	}

	dirEntries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal("Error Opening the dir", err)
	}

	for _, entry := range dirEntries {
		err := processImage(entry.Name(), inputDir, outputDir)
		if err == nil {
			log.Println("Successfully converted the file", entry.Name())
		}
	}

}
