package util

import (
	"errors"
	"fmt"

	"github.com/smallfish/simpleyaml"
)

// ConvertYamlListToStringList converts yaml list to list of strings
func ConvertYamlListToStringList(yamlList *simpleyaml.Yaml) ([]string, error) {
	yamlListArr, getArrayErr := yamlList.Array()
	if getArrayErr != nil {
		return nil, errors.New("Could not get Array for yaml list")
	}
	yamlListArrStrings := make([]string, len(yamlListArr))
	for index, yamlElement := range yamlListArr {
		yamlListArrStrings[index] = fmt.Sprint(yamlElement)
	}
	return yamlListArrStrings, nil
}
