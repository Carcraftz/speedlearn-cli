# speedlearn-cli

speedlearn-cli is a command-line tool that automates the process of downloading YouTube videos, converting them to MP3, and transcribing the audio content. This project was entirely built with Claude 3.5 Sonnet in under an hour.

## Features

- Download YouTube videos and convert them to MP3 format
- Transcribe audio files using the Groq API (Whisper model)
- Efficient file handling with caching to avoid redundant operations

## Prerequisites

- Go 1.x or higher
- FFmpeg (for video to audio conversion)
- Groq API key

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/speedlearn-cli.git
   cd speedlearn-cli
   ```

2. Install dependencies:
   ```
   make deps
   ```

3. Build the project:
   ```
   make build
   ```

## Configuration

Create a `.env` file in the project root and add your Groq API key:


## Usage

Run the tool with the following command:
./speedlearn-cli --filename <output_filename> --url <video_url>


The tool will:
1. Download the video and convert it to MP3
2. Transcribe the audio using the Groq API
3. Save the MP3 and transcription in the `output/<filename>` directory

## Project Structure

- `main.go`: Main application logic and CLI interface
- `convert.go`: Functions for video download and conversion to MP3
- `transcribe.go`: Functions for audio transcription using the Groq API
- `Makefile`: Build and management commands

## Development

- Run tests: `make test`
- Clean build artifacts: `make clean`
- Build for Linux: `make build-linux`

## Note

This entire project was conceptualized and implemented with the assistance of Claude 3.5 Sonnet, an AI language model, in under an hour. While functional, it may require further testing and refinement.