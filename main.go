package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Define command line flags
	filename := flag.String("filename", "", "Output filename (without extension)")
	videoURL := flag.String("url", "", "URL of the video file")
	flag.Parse()

	if *filename == "" || *videoURL == "" {
		fmt.Println("Both --filename and --url parameters are required")
		flag.Usage()
		os.Exit(1)
	}

	// Get API key from environment
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("GROQ_API_KEY environment variable is not set")
		os.Exit(1)
	}

	// Create output directory structure
	outputDir := filepath.Join("output", *filename)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Define output paths
	mp3Path := filepath.Join(outputDir, *filename+".mp3")
	txtPath := filepath.Join(outputDir, *filename+".txt")

	// Check if MP3 file exists
	if _, err := os.Stat(mp3Path); os.IsNotExist(err) {
		fmt.Printf("Downloading and converting video to MP3: %s\n", mp3Path)
		// Download and convert video to MP3
		err = convertVideoToMP3(*videoURL, mp3Path)
		if err != nil {
			fmt.Printf("Error converting video to MP3: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("MP3 file already exists: %s, skipping conversion\n", mp3Path)
	}

	// Check if transcription file exists
	if _, err := os.Stat(txtPath); os.IsNotExist(err) {
		fmt.Printf("Transcribing audio: %s\n", mp3Path)
		// Transcribe the audio
		transcription, err := transcribeAudio(mp3Path, apiKey)
		if err != nil {
			fmt.Printf("Error transcribing audio: %v\n", err)
			os.Exit(1)
		}

		// Save transcription to file
		fmt.Printf("Saving transcription to file: %s\n", txtPath)
		err = os.WriteFile(txtPath, []byte(transcription), 0644)
		if err != nil {
			fmt.Printf("Error saving transcription: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Transcription file already exists: %s\n", txtPath)
	}

	fmt.Printf("Processing complete:\nAudio: %s\nTranscription: %s\n", mp3Path, txtPath)
}
