package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

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
	// Define paths for minutia and notes
	minutiaPath := filepath.Join(outputDir, "minutia.md")
	notesPath := filepath.Join(outputDir, "notes.md")

	// Define note types and prompts
	type NoteConfig struct {
		Type   string
		Path   string
		Prompt string
	}

	noteConfigs := []NoteConfig{
		{
			Type:   "Minutia",
			Path:   minutiaPath,
			Prompt: "Give me any minutia or details that would be useful for an exam IN MARKDOWN. Be concise, don't explain things that would be known to someone that already understands the topic, just specific details. EG: OpenMP parallel for only parallelizes the outermost loop, not the inner loops.",
		},
		{
			Type:   "Detailed Notes",
			Path:   notesPath,
			Prompt: "Create very detailed notes from this lecture in markdown format. Include all important details, equations, and concepts. Add examples where appropriate. Please use markdown tables, code blocks, and other formatting where appropriate. Bold any important terms. Please include block quotes from the lecture.",
		},
	}

	// Get Anthropic API key from environment
	anthropicAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicAPIKey == "" {
		fmt.Println("ANTHROPIC_API_KEY environment variable is not set")
		os.Exit(1)
	}

	// Wrapper function to create notes
	createNotes := func(configs []NoteConfig, transcription string) error {
		for _, config := range configs {
			if _, err := os.Stat(config.Path); os.IsNotExist(err) {
				fmt.Printf("Creating %s from transcription: %s\n", config.Type, config.Path)

				notes, err := createNotesWithClaude(transcription, config.Prompt, anthropicAPIKey)
				if err != nil {
					return fmt.Errorf("error creating %s with Claude API: %v", config.Type, err)
				}

				err = ioutil.WriteFile(config.Path, []byte(notes), 0644)
				if err != nil {
					return fmt.Errorf("error saving %s: %v", config.Type, err)
				}
			} else {
				fmt.Printf("%s file already exists: %s\n", config.Type, config.Path)
			}
		}
		return nil
	}

	// Read transcription
	transcription, err := ioutil.ReadFile(txtPath)
	if err != nil {
		fmt.Printf("Error reading transcription file: %v\n", err)
		os.Exit(1)
	}

	// Create notes using the wrapper function
	err = createNotes(noteConfigs, string(transcription))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing complete:\nAudio: %s\nTranscription: %s\nMinutia: %s\n", mp3Path, txtPath, minutiaPath)
}
