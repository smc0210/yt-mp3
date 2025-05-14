# yt-mp3

`yt-mp3` is a program that extracts MP3 files from YouTube URLs and adds metadata and album artwork to them.

## Usage

### Running the Program

The program can be run using the `go run .` command. Use the following commands to perform specific tasks:

#### Download and Convert YouTube Video to MP3

```sh
go run . download
```

You will be prompted to enter the URL of the YouTube video. The program will download the video, convert it to an MP3 file, and save it in the `output` directory. Metadata and album artwork will also be added to the MP3 file.

#### Add Album Artwork to an Existing MP3 File

```sh
go run . addArtwork
```

You will be prompted to enter the path to the MP3 file and the path to the thumbnail image (JPG). The program will add the provided album artwork to the MP3 file.

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
