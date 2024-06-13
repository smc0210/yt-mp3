package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printUsage() {
	fmt.Println("\033[1;32mUsage:\033[0m")
	fmt.Println("  \033[1;36mgo run ytmp3.go download\033[0m    - Download and convert a YouTube video to MP3")
	fmt.Println("  \033[1;36mgo run ytmp3.go addArtwork\033[0m  - Add album artwork to an existing MP3 file")
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
	listFiles("./output")

	fmt.Print("Enter the path to the MP3 file: ")
	audioFile, _ := reader.ReadString('\n')
	audioFile = strings.TrimSpace(audioFile)
	audioFile = filepath.Join("./output", audioFile) // 경로를 절대 경로로 변환

	fmt.Print("Enter the path to the thumbnail image (JPG): ")
	thumbnailFile, _ := reader.ReadString('\n')
	thumbnailFile = strings.TrimSpace(thumbnailFile)
	thumbnailFile = filepath.Join("./output", thumbnailFile) // 경로를 절대 경로로 변환

	err := addArtwork(audioFile, thumbnailFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("앨범 아트워크가 성공적으로 추가되었습니다!")
}
