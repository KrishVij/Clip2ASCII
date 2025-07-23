package FFmpegutils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ExtarctFramesFromVideo(videoPath string) (framesPath string) {

	err := os.Mkdir("C:/Users/Krish Vij/Frames", 0750)

	framesPath = filepath.Join("C:/Users/Krish Vij","Frames")

	if err != nil {

		log.Fatalf("Error Occured while Creating the File: %v",err)
	}
	
	outputPattern := filepath.Join("C:/Users/Krish Vij/Frames", "%03d.png")
	
	cmd := exec.Command("ffmpeg", "-i", videoPath, outputPattern)
	output,err := cmd.Output()

	if err != nil {

		log.Fatalf("Error Occured while opening the ffmpeg Path: %v",err)
	}

	fmt.Println(string(output))

	return framesPath

}

// func StitchFramesToVideo(pathToFrames string) {

// 	inputFile,err := os.Open("C:/Users/Krish Vij/ASCII_Video")
// 	if err != nil {

// 		log.Fatalf("Error Occured while Creating the file: %v",err)
// 	}
// }