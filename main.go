package main

import "os"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "download", "d":
		handleDownload()
	case "addArtwork", "a":
		handleAddArtwork()
	default:
		printUsage()
	}
}
