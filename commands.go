package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func printUsage() {
	fmt.Println("\033[1;32mUsage:\033[0m")
	fmt.Println("  \033[1;36mgo run . download\033[0m    - Download and convert a YouTube video to MP3")
	fmt.Println("  \033[1;36mgo run . addArtwork\033[0m  - Add album artwork to an existing MP3 file")
	fmt.Println("  \033[1;36mgo run . video\033[0m       - Download a YouTube video (prompts for MP4 output)")
	fmt.Println("")
	fmt.Println("\033[1;32mCommands:\033[0m")
	fmt.Println("  \033[1;36mdownload (d)\033[0m     - Download and convert a YouTube video to MP3")
	fmt.Println("  \033[1;36maddArtwork (a)\033[0m   - Add album artwork to an existing MP3 file")
	fmt.Println("  \033[1;36mvideo (v)\033[0m        - Download a YouTube video (prompts for MP4 output)")
}

func handleDownload() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the URL of the YouTube video: ")
	videoURL, _ := reader.ReadString('\n')
	videoURL = strings.TrimSpace(videoURL)

	outputDir := "./output"
	videoTitle, err := getVideoTitle(videoURL)
	if err != nil {
		fmt.Printf("Error getting video title: %v\n", err)
		return
	}

	audioFile := filepath.Join(outputDir, videoTitle+".mp3")
	thumbnailFile := filepath.Join(outputDir, videoTitle+".webp")

	err = downloadAndConvert(videoURL, audioFile, thumbnailFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = addMetadata(audioFile, thumbnailFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("변환 및 태그 추가가 성공적으로 완료되었습니다!")
}

func handleAddArtwork() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\033[1;32mOutput 폴더에 있는 파일 목록:\033[0m")
	mp3Files, imageFiles := listFilesWithIndexes("./output")

	fmt.Print("Enter the index of the MP3 file: ")
	audioIndexStr, _ := reader.ReadString('\n')
	audioIndex, err := strconv.Atoi(strings.TrimSpace(audioIndexStr))
	if err != nil || audioIndex < 0 || audioIndex >= len(mp3Files) {
		fmt.Println("Invalid index.")
		return
	}
	audioFile := filepath.Join("./output", mp3Files[audioIndex])

	fmt.Print("Enter the index of the thumbnail image (JPG): ")
	thumbnailIndexStr, _ := reader.ReadString('\n')
	thumbnailIndex, err := strconv.Atoi(strings.TrimSpace(thumbnailIndexStr))
	if err != nil || thumbnailIndex < 0 || thumbnailIndex >= len(imageFiles) {
		fmt.Println("Invalid index.")
		return
	}
	thumbnailFile := filepath.Join("./output", imageFiles[thumbnailIndex])

	err = addArtwork(audioFile, thumbnailFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("앨범 아트워크가 성공적으로 추가되었습니다!")
}

func handleVideoDownload() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the URL of the YouTube video: ")
	videoURL, _ := reader.ReadString('\n')
	videoURL = strings.TrimSpace(videoURL)

	fmt.Println("\033[1;32mSelect video quality:\033[0m")
	fmt.Println("  [1] 720p")
	fmt.Println("  [2] 1080p")
	fmt.Println("  [3] Best available")
	fmt.Print("Enter quality option (1-3): ")
	qualityChoiceStr, _ := reader.ReadString('\n')
	qualityChoice, err := strconv.Atoi(strings.TrimSpace(qualityChoiceStr))
	if err != nil || qualityChoice < 1 || qualityChoice > 3 {
		fmt.Println("Invalid quality option. Defaulting to Best available.")
		qualityChoice = 3 // Default to best if input is invalid
	}

	var qualityOption string
	switch qualityChoice {
	case 1:
		qualityOption = "720p"
	case 2:
		qualityOption = "1080p"
	case 3:
		qualityOption = "best"
	default:
		qualityOption = "best"
	}

	fmt.Print("Force MP4 output? (y/N): ")
	forceMp4Str, _ := reader.ReadString('\n')
	forceMp4 := strings.TrimSpace(strings.ToLower(forceMp4Str)) == "y"

	outputDir := "./output"
	videoTitle, err := getVideoTitle(videoURL)
	if err != nil {
		fmt.Printf("Error getting video title: %v\n", err)
		return
	}

	// Determine the expected extension based on forceMp4, mainly for the success message
	// yt-dlp will handle the actual extension if not forcing mp4 for "best"
	// expectedExtension := ".mp4" // Default to .mp4 for messages if forcing or specific quality
	// if qualityOption == "best" && !forceMp4 {
	// 	// If best quality and not forcing MP4, we can't be certain of the extension.
	// 	// For the success message, we'll still say .mp4 but acknowledge yt-dlp might choose differently.
	// 	// Or, we can adjust the downloadVideo function to return the actual filename.
	// 	// For now, let's keep it simple and the message might be slightly off for "best" non-forced.
	// 	// A better approach would be for downloadVideo to return the final filename.
	// }

	// videoFile := filepath.Join(outputDir, videoTitle+expectedExtension) // This line is no longer needed

	err = downloadVideo(videoURL, outputDir, videoTitle, qualityOption, forceMp4)
	if err != nil {
		fmt.Printf("Error downloading video: %v\n", err)
		return
	}

	// To get the *actual* final filename, we might need to list dir or have downloadVideo return it.
	// For now, the success message uses the 'expected' name.
	fmt.Printf("Video download process completed. Check %s for the file (extension may vary if not forcing MP4 for 'best' quality).\n", outputDir)
}
