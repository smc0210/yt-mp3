package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func getVideoTitle(videoURL string) (string, error) {
	cmd := exec.Command("yt-dlp", "--get-filename", "-o", "%(title)s", videoURL)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	title := out.String()
	title = strings.TrimSpace(title)

	// Normalize Unicode to ensure proper handling of combining characters
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}
	title = strings.Map(func(r rune) rune {
		if isMn(r) {
			return -1
		}
		return r
	}, norm.NFC.String(title))

	title = sanitizeFilename(title)
	if len(title) > 50 {
		title = title[:50]
	}
	return title, nil
}

func downloadAndConvert(videoURL, audioFile, thumbnailFile string) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(audioFile)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	// Sanitize file names
	baseAudioFileName := sanitizeFilename(strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile)))
	baseThumbnailFileName := sanitizeFilename(strings.TrimSuffix(filepath.Base(thumbnailFile), filepath.Ext(thumbnailFile)))

	// Construct full paths
	sanitizedAudioFile := filepath.Join(outputDir, baseAudioFileName+".mp3")
	sanitizedThumbnailFile := filepath.Join(outputDir, baseThumbnailFileName)

	// Download video and convert to MP3 using yt-dlp and ffmpeg
	fmt.Println("Downloading and converting video to MP3...")
	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "mp3", "--audio-quality", "0", "--embed-metadata", "--parse-metadata", "title:%(meta_title)s", "--output", sanitizedAudioFile, videoURL)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error during download and convert: %v\n", err)
		return err
	}

	// Download thumbnail
	fmt.Println("Downloading thumbnail...")
	cmd = exec.Command("yt-dlp", "--write-thumbnail", "--skip-download", "--convert-thumbnails", "jpg", "--output", sanitizedThumbnailFile+".%(ext)s", videoURL)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err = cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error during thumbnail download: %v\n", err)
		return err
	}

	// Find the converted jpg file
	fullThumbnailFilePath := sanitizedThumbnailFile + ".jpg"
	if _, err := os.Stat(fullThumbnailFilePath); err != nil {
		return fmt.Errorf("Converted thumbnail not found: %v", err)
	}

	return nil
}
