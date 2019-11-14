package util

import (
	"fmt"

	"github.com/smallfish/simpleyaml"
)

// ConvertYamlListToStringList converts yaml list to list of strings
func ConvertYamlListToStringList(yamlList *simpleyaml.Yaml) ([]string, error) {
	yamlListArr, getArrayErr := yamlList.Array()
	if getArrayErr != nil {
		return nil, fmt.Errorf(
			"Could not build array from yaml list\n%s",
			getArrayErr,
		)
	}
	yamlListArrStrings := make([]string, len(yamlListArr))
	for index, yamlElement := range yamlListArr {
		yamlListArrStrings[index] = fmt.Sprint(yamlElement)
	}
	return yamlListArrStrings, nil
}
