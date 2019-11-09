package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
)

// ReadFileFromPathOrURL read file from either os filesystem
// or from a URL
func ReadFileFromPathOrURL(source string) ([]byte, error) {

	if fileExists(source) {
		sourceAbsPath, absPathErr := filepath.Abs(source)
		if absPathErr != nil {
			return nil, absPathErr
		}
		return ioutil.ReadFile(sourceAbsPath)
	} else if validURL(source) {
		fileBody, downloadFileErr := downloadFile(source)
		if downloadFileErr != nil {
			return nil, downloadFileErr
		}
		return fileBody, nil
	}

	return nil, errors.New("Could not find runes.yaml file")
}

func getWorkingDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func validURL(source string) bool {
	return govalidator.IsURL(source)
}

func fileExists(source string) bool {
	fileInfo, fileStatErr := os.Stat(source)
	if os.IsNotExist(fileStatErr) {
		return false
	}
	return !fileInfo.IsDir()
}

func downloadFile(url string) ([]byte, error) {
	httpResp, httpErr := http.Get(url)
	if httpErr != nil {
		return nil, httpErr
	}

	var data bytes.Buffer
	_, copyErr := io.Copy(&data, httpResp.Body)
	if copyErr != nil {
		return nil, copyErr
	}
	return data.Bytes(), nil
}
