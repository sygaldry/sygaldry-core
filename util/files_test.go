package util

import (
	"testing"
)

// TestGetWorkingDir tests util.getWorkingDir
func TestGetWorkingDir(t *testing.T) {
	_, getWorkingDirErr := getWorkingDir()
	if getWorkingDirErr != nil {
		t.Error("Tried to get current working directory, got", getWorkingDirErr)
	}
}

func TestValidURL(t *testing.T) {
	knownValidURL := "https://golang.org/"
	knownInvalidURL := "relative/path/to/golang"

	knownValidURLIsValid := validURL(knownValidURL)
	if !knownValidURLIsValid {
		t.Error("Tried to validate known valid URL, expected true, got", knownValidURLIsValid)
	}

	knownInvalidURLIsValid := validURL(knownInvalidURL)
	if knownInvalidURLIsValid {
		t.Error("Tried to validate known invalid URL, expected false, got", knownInvalidURLIsValid)
	}
}

func TestDownloadFile(t *testing.T) {
	downloadURL := "https://golang.org/lib/godoc/images/footer-gopher.jpg"
	imageBody, downloadError := downloadFile(downloadURL)
	if downloadError != nil {
		t.Error("Tried to download gopher from URL, but got ", downloadError)
	}
	if len(imageBody) == 0 {
		t.Error("Tried to download image, but got this many bytes", len(imageBody))
	}
}

func TestReadFileFromPathOrURL(t *testing.T) {
	_, runesFileReadErr := ReadFileFromPathOrURL(
		"test_resources/test-runes.yaml",
	)
	if runesFileReadErr != nil {
		t.Error("Tried to read test runes.yaml file, got", runesFileReadErr)
	}
}

func TestFileExists(t *testing.T) {
	knownExistingFile := "files.go"
	knownNotExistingFile := "files.not.a.real.file.go"
	if !fileExists(knownExistingFile) {
		t.Error("Tried to check existence of existing file, expected true, got", knownExistingFile)
	}
	if fileExists(knownNotExistingFile) {
		t.Error("Tried to check existence of not existing file, expected false, got", knownNotExistingFile)
	}
}
