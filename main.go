package main

import (

	"fmt"
	"log"
	"os"

	"path/filepath"

	FFmpegutils "github.com/KrishVij/clip2ASCII/FFmpeg_Utils"
	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

func main() {

	const videoPath = "C:/Users/Krish Vij/Downloads/vecteezy_visiting-a-destination-city-on-holiday_34634255.mov"
	const outputPATH = "C:/Users/Krish Vij/output.mp4"

	err := os.MkdirAll("C:/Users/Krish Vij/ASCII_Frames", 0750)
	if err != nil {

		log.Fatalf("Error creating ASCII_Frames directory: %v", err)
	}
	
	Count := 1

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
		Frame_Processing.SaveImage(asciiImage, Count)
		if err != nil {

			log.Fatalf("Error occured while saving the image: %v", err)
		}

		Count++

		fmt.Printf("Frame: %d processed successfully\n",Count - 1)

	}

	log.Println("ASCII frames folder generated successfully!")

	FFmpegutils.StitchFramesToVideo(outputPATH)

	fmt.Println("ASCII image Generated successfully")

}
