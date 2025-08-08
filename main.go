package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"path/filepath"

	FFmpegutils "github.com/KrishVij/clip2ASCII/FFmpeg_Utils"
	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

func main() {

	const videoPath = "C:/Users/Krish Vij/Downloads/testvid3.mp4"
	const outputPATH = "C:/Users/Krish Vij/output3.mp4"

	// mpp := Frame_Processing.Calculate_Ink_Required_For_Drawing_ASCII_Chars()
	// Frame_Processing.Generate_ASCII_Lookup_Table(mpp)

	err := os.MkdirAll("C:/Users/Krish Vij/ASCII_Frames", 0750)
	if err != nil {

		log.Fatalf("Error creating ASCII_Frames directory: %v", err)
	}

	// Count := 1
	
	result := FFmpegutils.ExtarctFramesFromVideo(videoPath)

	directory, err := os.ReadDir(result)
	if err != nil {

		log.Fatalf("Error Occured While Reading Contents of The Frames Folder: %v", err)
	}

	var wg sync.WaitGroup

	for i, item := range directory {

		wg.Add(1)

		go func(count int, filename string) {

			defer wg.Done()

			framePath := filepath.Join(result, filename)

		    Image, err := Frame_Processing.LoadAndResizeImage(framePath)

		   if err != nil {

			log.Fatalf("ERROR OCCURED WHILE LOADING IMAGE: %v", err)
		   }

		    Pixels, rgbaValues := Frame_Processing.ExtractPixelData(Image)
		    asciiImage, err := Frame_Processing.RenderAsciiImage(Pixels, rgbaValues)
		    if err != nil {

			log.Fatalf("ERROR OCCURED WHILE CONVERTING TO ASCII: %v", err)

		    }
		    err = Frame_Processing.SaveImage(asciiImage, count)
		    if err != nil {

			log.Fatalf("Error occured while saving the image: %v", err)
		    }
			
			fmt.Printf("Frame: %d processed successfully\n", count-1)
		}(i + 1, item.Name())

	}

	wg.Wait()

	log.Println("ASCII frames folder generated successfully!")

	FFmpegutils.StitchFramesToVideo(outputPATH)

	fmt.Println("ASCII Video Generated successfully")

}
