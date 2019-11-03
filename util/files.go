package util

import (
	"errors"
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
	workingDir, getWorkingDirErr := getWorkingDir()
	if getWorkingDirErr != nil {
		return nil, getWorkingDirErr
	}

	var runesYamlPath string

	if fileExists(source) {
		sourceAbsPath, absPathErr := filepath.Abs(source)
		if absPathErr != nil {
			return nil, absPathErr
		}
		runesYamlPath = sourceAbsPath
	} else if validURL(source) {
		runesYamlPath = fmt.Sprintf("%s/runes.yaml", workingDir)
		if downloadFileErr := downloadFile(runesYamlPath, source); downloadFileErr != nil {
			return nil, downloadFileErr
		}
	} else {
		return nil, errors.New("Could not find runes.yaml file")
	}

	return ioutil.ReadFile(runesYamlPath)
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

func downloadFile(filepath string, url string) error {

	file, osErr := os.Create(filepath)
	if osErr != nil {
		return osErr
	}
	defer file.Close()

	httpResp, httpErr := http.Get(url)
	if httpErr != nil {
		return httpErr
	}
	defer httpResp.Body.Close()

	_, copyErr := io.Copy(file, httpResp.Body)
	return copyErr
}
