package FFmpegutils

import (
	
	"fmt"
	"log"
	"bytes"
	"strings"
	"strconv"
	"os"
	"os/exec"
	"path/filepath"
)

func ExtractFramesFromVideo(videoPath string) (framesPath string) {

	user_home_directory, err := os.UserHomeDir()
	if err != nil {

		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}
	
	path_to_Frames := filepath.Join(user_home_directory, "Frames")
	err = os.Mkdir(path_to_Frames, 0750)
	if err != nil {

		log.Fatalf("Error Occured while Creating the Directory: %v", err)
	}

	framesPath = path_to_Frames
	outputPattern := filepath.Join(path_to_Frames, "%03d.png")

	cmd := exec.Command("ffmpeg", "-i", videoPath, outputPattern)
	output, err := cmd.CombinedOutput()

	if err != nil {

		log.Fatalf("Error Occured while opening the ffmpeg Path: %v", err)
	}

	fmt.Println(string(output))

	return framesPath

}

func StitchFramesToVideo(outputPATH string) {

	user_home_directory, err := os.UserHomeDir()
	if err != nil {

		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}
	
	path_to_ASCII_FRAMES := filepath.Join(user_home_directory, "ASCII_FRAMES")
	err = os.Chdir(path_to_ASCII_FRAMES)
	if err != nil {
		
		log.Fatalf("Error changing directory: %v", err)
	}

	cmd := exec.Command("ffmpeg", "-framerate", "30", "-i", "ASCII_Frames%03d.png", outputPATH)
	output, err := cmd.CombinedOutput()

	if err != nil {

		log.Fatalf("FFmpeg error: %v\nOutput: %s", err, string(output))
	}

}

func Extract_Thumbnail(videoPath string) (thumbnail_file_path string) {

	user_home_directory, err := os.UserHomeDir()
	if err != nil {

		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}
	path_to_thumbnail_directory := filepath.Join(user_home_directory, "thumbnail")
	err = os.Mkdir(path_to_thumbnail_directory, 0750)
	if err != nil {

		log.Fatalf("Error Occured while Creating the Directory: %v", err)

	}
	
	thumbnail_file_path = filepath.Join(path_to_thumbnail_directory, "thumbnail.png")
	
	cmd  := exec.Command("ffmpeg","-i",videoPath,"-ss","0", "-vframes", "1", thumbnail_file_path)

	cmd.Stdout = os.Stdout  // Capture stdout.
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {

		log.Fatalf("FFmpeg Error : %v",err)
	}

	return thumbnail_file_path
}

func Check_Duration(videoPath string) bool {

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error occurred while checking video duration: %v", err)
	}

	durationStr := strings.TrimSpace(string(bytes.Trim(output, "\n")))
	if durationStr == "" {
		log.Fatalf("No duration found in ffprobe output")
	}

	val, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		log.Fatalf("Couldn't convert duration to float64: %v", err)
	}

	return val <= 30
}
