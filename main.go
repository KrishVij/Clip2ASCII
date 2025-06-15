package main

import (
	"log"

	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

func main() {

	Image, err := Frame_Processing.ImageResizeAndLoad()

	if err != nil {

		log.Fatalf("ERROR OCCURED WHILE LOADING IMAGE: %v", err)
	}

	Pixels,rgbaValues := Frame_Processing.ProcessImageForAscii(Image)
	asciiImage,err := Frame_Processing.GrayScaleToAscii(Pixels,rgbaValues)
	if err != nil {

		log.Fatalf("ERROR OCCURED WHILE CONVERTING TO ASCII: %v",err)
		
	}
	Frame_Processing.SaveImage(asciiImage)
	if err != nil {

		log.Fatalf("Error occured while saving the image: %v",err)
	}

	log.Println("ASCII image generated successfully!")

}
