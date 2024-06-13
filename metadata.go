package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogem/id3v2"
)

func addMetadata(audioFile, thumbnailFile string) error {
	baseAudioFileName := sanitizeFilename(strings.TrimSuffix(filepath.Base(audioFile), filepath.Ext(audioFile)))
	fullAudioFilePath := filepath.Join(filepath.Dir(audioFile), baseAudioFileName+".mp3")
	tag, err := id3v2.Open(fullAudioFilePath, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	baseThumbnailFileName := sanitizeFilename(strings.TrimSuffix(filepath.Base(thumbnailFile), filepath.Ext(thumbnailFile)))
	fullThumbnailFilePath := filepath.Join(filepath.Dir(thumbnailFile), baseThumbnailFileName+".jpg")

	if thumbnailData, err := os.ReadFile(fullThumbnailFilePath); err == nil {
		tag.AddAttachedPicture(id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Description: "Thumbnail",
			Picture:     thumbnailData,
		})
	} else {
		fmt.Printf("Error adding thumbnail: %v\n", err)
		return fmt.Errorf("Error adding thumbnail: %v", err)
	}

	if err := tag.Save(); err != nil {
		return err
	}

	fmt.Println("MP3 파일에 추가된 메타데이터:")
	for k, frames := range tag.AllFrames() {
		for _, frame := range frames {
			switch f := frame.(type) {
			case id3v2.TextFrame:
				fmt.Printf("%s: %s\n", k, f.Text)
			case id3v2.PictureFrame:
				fmt.Printf("%s: [Picture]\n", k)
			default:
				if txxFrame, ok := frame.(id3v2.UserDefinedTextFrame); ok {
					fmt.Printf("%s: %s (%s)\n", k, txxFrame.Description, txxFrame.Value)
				}
			}
		}
	}

	return nil
}

func addArtwork(audioFile, thumbnailFile string) error {
	tag, err := id3v2.Open(audioFile, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	if thumbnailData, err := os.ReadFile(thumbnailFile); err == nil {
		tag.AddAttachedPicture(id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Description: "Thumbnail",
			Picture:     thumbnailData,
		})
	} else {
		fmt.Printf("Error adding thumbnail: %v\n", err)
		return fmt.Errorf("Error adding thumbnail: %v", err)
	}

	if err := tag.Save(); err != nil {
		return err
	}

	fmt.Println("앨범 아트워크가 성공적으로 추가되었습니다!")
	return nil
}
