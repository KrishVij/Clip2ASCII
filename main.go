package main

import (
	"log"
	"os"

	"path/filepath"

	FFmpegutils "github.com/KrishVij/clip2ASCII/FFmpeg_Utils"
	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

func main() {

	const videoPath = "C:/Users/Krish Vij/Downloads/vecteezy_visiting-a-destination-city-on-holiday_34634255.mov"

	result := FFmpegutils.ExtarctFramesFromVideo(videoPath)

	directory, err := os.ReadDir(result)
	if err != nil {

		log.Fatalf("Error Occured While Reading Contents of The Frames Folder: %v", err)
	}

	for _, item := range directory {

		framePath := filepath.Join(result, item.Name())

		Image, err := Frame_Processing.LoadAndResizeImage(framePath)

		if err != nil {

			log.Fatalf("ERROR OCCURED WHILE LOADING IMAGE: %v", err)
		}

		Pixels, rgbaValues := Frame_Processing.ExtractPixelData(Image)
		asciiImage, err := Frame_Processing.RenderAsciiImage(Pixels, rgbaValues)
		if err != nil {

			log.Fatalf("ERROR OCCURED WHILE CONVERTING TO ASCII: %v", err)

		}
		Frame_Processing.SaveImage(asciiImage)
		if err != nil {

			log.Fatalf("Error occured while saving the image: %v", err)
		}

		log.Println("ASCII image generated successfully!")

	}

	/*
	
	Image, err := Frame_Processing.LoadAndResizeImage()

	if err != nil {

		log.Fatalf("ERROR OCCURED WHILE LOADING IMAGE: %v", err)
	}

	Pixels, rgbaValues := Frame_Processing.ExtractPixelData(Image)
	asciiImage, err := Frame_Processing.RenderAsciiImage(Pixels, rgbaValues)
	if err != nil {

		log.Fatalf("ERROR OCCURED WHILE CONVERTING TO ASCII: %v", err)

	}
	Frame_Processing.SaveImage(asciiImage)
	if err != nil {

		log.Fatalf("Error occured while saving the image: %v", err)
	}

	log.Println("ASCII image generated successfully!")
	*/

}
