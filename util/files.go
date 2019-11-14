package util

import (
	"bytes"
	"fmt"
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
			return nil, fmt.Errorf(
				"Could not download yaml file at %s\n%s",
				source,
				downloadFileErr,
			)
		}
		return fileBody, nil
	}

	return nil, fmt.Errorf(
		"Could not find yaml file at %s",
		source,
	)
}

func getWorkingDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf(
			"Could not determine current working directory\n%s",
			err,
		)
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
		return nil, fmt.Errorf(
			"Could not GET url %s\n%s",
			url,
			httpErr,
		)
	}

	var data bytes.Buffer
	_, copyErr := io.Copy(&data, httpResp.Body)
	if copyErr != nil {
		return nil, fmt.Errorf(
			"Could not copy GET response from %s to buffer\n%s",
			url,
			copyErr,
		)
	}
	return data.Bytes(), nil
}
