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
	fmt.Println("")
	fmt.Println("\033[1;32mCommands:\033[0m")
	fmt.Println("  \033[1;36mdownload (d)\033[0m     - Download and convert a YouTube video to MP3")
	fmt.Println("  \033[1;36maddArtwork (a)\033[0m   - Add album artwork to an existing MP3 file")
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
