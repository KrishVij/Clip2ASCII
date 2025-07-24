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

	framesPath = filepath.Join("C:/Users/Krish Vij", "Frames")

	if err != nil {

		log.Fatalf("Error Occured while Creating the File: %v", err)
	}

	outputPattern := filepath.Join("C:/Users/Krish Vij/Frames", "%03d.png")

	cmd := exec.Command("ffmpeg", "-i", videoPath, outputPattern)
	output, err := cmd.Output()

	if err != nil {

		log.Fatalf("Error Occured while opening the ffmpeg Path: %v", err)
	}

	fmt.Println(string(output))

	return framesPath

}

func StitchFramesToVideo(outputPATH string) {

	err := os.Chdir("C:/Users/Krish Vij/ASCII_Frames")
	
	if err != nil {
		
		log.Fatalf("Error changing directory: %v", err)
	}

	cmd := exec.Command("ffmpeg", "-framerate", "30", "-i", "ASCII_Frames%03d.png", outputPATH)
	output, err := cmd.CombinedOutput()

	if err != nil {

		log.Fatalf("FFmpeg error: %v\nOutput: %s", err, string(output))
	}

}
