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

func downloadVideo(videoURL, outputDir, videoTitle, qualityOption string, forceMp4 bool) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Base filename for yt-dlp (without extension, yt-dlp will add it or use container's default)
	// We will construct the output template for yt-dlp to place it directly in the outputDir with videoTitle
	outputPathWithTitle := filepath.Join(outputDir, videoTitle)

	var formatString string
	switch qualityOption {
	case "720p":
		formatString = "bestvideo[height<=720]+bestaudio/best[height<=720]"
	case "1080p":
		formatString = "bestvideo[height<=1080]+bestaudio/best[height<=1080]"
	case "best":
		formatString = "bestvideo+bestaudio/best"
	default:
		formatString = "bestvideo+bestaudio/best" // Default to best
	}

	fmt.Printf("Downloading video with quality '%s'...\n", qualityOption)
	
	args := []string{"-f", formatString}
	if forceMp4 {
		args = append(args, "--merge-output-format", "mp4")
		// If forcing mp4, we explicitly set the output filename with .mp4 extension
		// yt-dlp's -o flag will use this complete path.
		outputPathWithTitleAndExt := outputPathWithTitle + ".mp4"
		args = append(args, "-o", outputPathWithTitleAndExt)
	} else {
		// If not forcing mp4, let yt-dlp decide the extension. 
		// We provide the path and base name, yt-dlp appends appropriate extension.
		// For example, 'output/My Video' becomes 'output/My Video.mkv' or 'output/My Video.mp4'
		args = append(args, "-o", outputPathWithTitle+".%(ext)s")
	}
	args = append(args, videoURL)

	cmd := exec.Command("yt-dlp", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start video download command: %w", err)
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
			fmt.Println(scanner.Text()) // yt-dlp progress often goes to stderr
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("video download command failed: %w", err)
	}

	// The success message in commands.go is now more generic about the output directory.
	return nil
}
