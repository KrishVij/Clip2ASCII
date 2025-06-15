package main

import (

	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

func main() {

	Image := Frame_Processing.LoadImage()

	if Image != nil {

		defer Image.Close()
	}

	Pixels := Frame_Processing.GetBrightnessMatrix(Image)
	asciiImage := Frame_Processing.GrayScaleToAscii(Pixels)
	Frame_Processing.SaveImage(asciiImage)

}
