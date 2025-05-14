# yt-mp3

`yt-mp3` is a program that extracts MP3 files from YouTube URLs and adds metadata and album artwork to them.

## Usage

### Running the Program

The program is run from the terminal using `go run .` followed by a command.
If no command is provided or an invalid command is used, a list of available commands will be shown.

Here's what each command does:

*   **Download and Convert YouTube Video to MP3**
    *   Command: `go run . download` (or `d`)
    *   Prompts you to enter the URL of the YouTube video.
    *   Downloads the audio, converts it to an MP3 file, adds metadata (like title and artwork from YouTube if available), and saves it in the `output` directory.

*   **Download YouTube Video as MP4 (or other best format)**
    *   Command: `go run . video` (or `v`)
    *   Prompts you to enter the URL of the YouTube video.
    *   Asks you to select the video quality (e.g., 720p, 1080p, best available).
    *   Asks if you want to force the output to be an MP4 file.
    *   Downloads the video to the `output` directory.
        *   *Note:* If you select "best available" quality and don't force MP4 output, `yt-dlp` might save the file in a different container format (like MKV) if that's better for the chosen video/audio streams.

*   **Add Album Artwork to an Existing MP3 File**
    *   Command: `go run . addArtwork` (or `a`)
    *   Lists MP3 and JPG image files currently in the `output` directory.
    *   Prompts you to enter the index number for the MP3 file and for the JPG image.
    *   Embeds the selected image as album artwork into the MP3 file.

## Project Structure

The project is structured as follows:

```
yt-mp3/
├── go.mod
├── main.go
├── commands.go
├── download.go
├── metadata.go
└── utils.go
```

- `main.go`: Entry point of the program and basic command handling logic.
- `commands.go`: Functions that handle user commands.
- `download.go`: Logic for downloading and converting YouTube videos.
- `metadata.go`: Logic for adding metadata and album artwork to MP3 files.
- `utils.go`: Utility functions.

## Requirements

Make sure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.16 or later)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
  - **Installation:**
    - **macOS (using [Homebrew](https://brew.sh/)):**
      ```sh
      brew install yt-dlp
      ```
    - **Other OS / Manual Installation:** Please refer to the [official yt-dlp installation guide](https://github.com/yt-dlp/yt-dlp#installation).
- [ffmpeg](https://ffmpeg.org/)

## Installation

1. Clone the repository:

```sh
git clone https://github.com/yourusername/yt-mp3.git
cd yt-mp3
```

2. Initialize Go modules:

```sh
go mod tidy
```

3. Run the program:

```sh
go run .
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
