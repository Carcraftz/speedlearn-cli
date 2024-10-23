package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func createNotesWithClaude(transcription, prompt, apiKey string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":      "claude-3-5-sonnet-latest",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "user", "content": prompt + "\n\n" + transcription},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("unexpected response format: 'content' is not a non-empty slice. Full response: %s", string(body))
	}

	firstContent, ok := content[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected response format: first element of 'content' is not a map. Full response: %s", string(body))
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: 'text' key not found or not a string. Full response: %s", string(body))
	}

	return text, nil
}
