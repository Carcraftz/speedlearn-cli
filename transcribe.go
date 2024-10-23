package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func transcribeAudio(audioPath, apiKey string) (string, error) {
	// Open the audio file
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}
	defer file.Close()

	// Prepare the multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file to form: %v", err)
	}

	// Add other form fields
	writer.WriteField("model", "whisper-large-v3-turbo")
	writer.WriteField("language", "en")
	writer.WriteField("response_format", "verbose_json")

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result TranscriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return strings.TrimSpace(result.Text), nil
}
