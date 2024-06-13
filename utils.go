package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func listFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	mp3Files := []string{}
	imageFiles := []string{}

	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".mp3") {
			mp3Files = append(mp3Files, name)
		} else if strings.HasSuffix(name, ".jpg") {
			imageFiles = append(imageFiles, name)
		}
	}

	fmt.Println("\033[1;36mMP3 Files:\033[0m")
	for _, file := range mp3Files {
		fmt.Printf("  - %s\n", file)
	}

	fmt.Println("\033[1;36mImage Files:\033[0m")
	for _, file := range imageFiles {
		fmt.Printf("  - %s\n", file)
	}
}

func sanitizeFilename(filename string) string {
	// 비표준 유니코드 문자 및 제어 문자 제거
	sanitized := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || r == unicode.ReplacementChar || (r >= 0xD800 && r <= 0xDFFF) || r > 0x10FFFF {
			return -1
		}
		return r
	}, filename)

	// 특수 문자 제거
	re := regexp.MustCompile(`[^\w\s-가-힣]`)
	sanitized = re.ReplaceAllString(sanitized, "")

	// 공백과 하이픈을 제외한 모든 것을 제거
	sanitized = strings.ReplaceAll(sanitized, " ", "_")

	// 파일명 길이 제한
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return strings.TrimSpace(sanitized)
}
