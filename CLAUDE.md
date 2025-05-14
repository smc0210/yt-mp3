# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- Run application: `go run .`
- Build executable: `go build -o yt-mp3`
- Format code: `go fmt ./...`
- Check dependencies: `go mod tidy`

## Code Style Guidelines
- Follow standard Go conventions (camelCase for variables/functions)
- Error handling: Return errors to caller, check after operations
- File permissions: Use octal notation (0o755)
- Maintain Korean text for user-facing messages
- Keep filenames lowercase and function-specific
- Maximum filename length: 50 characters
- Sanitize user inputs especially for filenames
- Organize imports alphabetically (standard lib first, then external)
- Use appropriate Unicode normalization for Korean text
- Maximum line length ~80 characters
- Use early returns for error conditions

## Dependencies
- Requires external tools: yt-dlp, ffmpeg
- Go dependencies: github.com/bogem/id3v2, golang.org/x/text