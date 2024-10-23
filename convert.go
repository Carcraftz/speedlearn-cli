package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func convertVideoToMP3(videoURL, outputPath string) error {
	// Ensure output directory exists
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Build ffmpeg command
	cmd := exec.Command("ffmpeg",
		"-i", videoURL,
		"-b:a", "32k", // Reduced bitrate from 64k to 32k for further compression
		"-vn",
		"-acodec", "libmp3lame", // Explicitly specify the MP3 codec
		"-q:a", "9", // Use variable bitrate (VBR) encoding with quality level 9 (lowest quality, highest compression)
		outputPath,
	)

	// Capture both stdout and stderr and display them in real-time
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v", err)
	}

	fmt.Println("Video conversion completed successfully.")
	return nil
}
