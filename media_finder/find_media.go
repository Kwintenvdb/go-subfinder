package mediafinder

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	videoFileExtensions = []string{
		".webm",
		".mkv",
		".flv",
		".vob",
		".ogv",
		".ogg",
		".rrc",
		".gif",
		".mng",
		".mov",
		".avi",
		".qt",
		".wmv",
		".yuv",
		".rm",
		".asf",
		".amv",
		".mp4",
		".m4p",
		".m4v",
		".mpg",
		".mp2",
		".mpeg",
		".mpe",
		".mpv",
		".m4v",
		".svi",
		".3gp",
		".3g2",
		".mxf",
		".roq",
		".nsv",
		".flv",
		".f4v",
		".f4p",
		".f4a",
		".f4b",
	}
)

func FindVideoFileInCurrentDir() (fileName string, err error) {
	files := getFilesInCurrentDir()

	videoFiles := make([]fs.FileInfo, 0)
	for _, f := range files {
		if isVideoFile(f) {
			videoFiles = append(videoFiles, f)
		}
	}

	numVideoFiles := len(videoFiles)

	if numVideoFiles == 0 {
		return "", errors.New("no video file found in current dir")
	}

	if numVideoFiles > 1 {
		return "", fmt.Errorf("found more than 1 video file. Found %d", numVideoFiles)
	}

	return videoFiles[0].Name(), nil
}

func getFilesInCurrentDir() []fs.FileInfo {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func isVideoFile(file fs.FileInfo) bool {
	fileExt := filepath.Ext(file.Name())
	for _, ext := range videoFileExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}
