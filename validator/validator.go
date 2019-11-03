package validator

import (
	"errors"
	"fmt"

	"github.com/smallfish/simpleyaml"
	"github.com/sygaldry/sygaldry-core/runes"
	"github.com/sygaldry/sygaldry-core/util"
	"gopkg.in/yaml.v2"
)

// GetValidStage validates a rune yaml file
// either a path on the local filesystem or a url
// is passed in
func GetValidStage(source string, stage string) ([]runes.Rune, error) {
	runeYaml, newYamlErr := yamlSource(source)
	if newYamlErr != nil {
		return nil, newYamlErr
	}

	runeDefinitions, getRuneDefinitionsErr := getRuneDefinitions(runeYaml)
	if getRuneDefinitionsErr != nil {
		return nil, getRuneDefinitionsErr
	}
	runeConfigs, getRuneConfigsErr := getRunesConfigsForStage(runeYaml, stage)
	if getRuneConfigsErr != nil {
		return nil, getRuneConfigsErr
	}

	runeConfigsArray, arrayErr := runeConfigs.Array()

	if arrayErr != nil {
		return nil, arrayErr
	}

	for _, element := range runeConfigsArray {

	}

	return []runes.Rune{}, nil
}

func validateRuneConfig(runeConfig *simpleyaml.Yaml, runeDefinitions *simpleyaml.Yaml) (bool, error) {

}

func getRunesConfigsForStage(runeYaml *simpleyaml.Yaml, stage string) (*simpleyaml.Yaml, error) {
	stageRunes := runeYaml.Get("stages").Get(stage)

	if numRunes, _ := stageRunes.GetArraySize(); numRunes == 0 {
		return nil, fmt.Errorf(
			"Could not find any Runes for stage %s",
			stage,
		)
	}
	return stageRunes, nil
}

func getRuneDefinitions(runeYaml *simpleyaml.Yaml) (*simpleyaml.Yaml, error) {
	definitionURLs := runeYaml.Get("definitions")
	definitionURLsArray, getDefinitionURLsErr := definitionURLs.Array()
	if getDefinitionURLsErr != nil {
		return nil, errors.New("could not find rune definitions")
	}

	definitions := map[interface{}]interface{}{}

	for _, definitionURL := range definitionURLsArray {
		definitionURLString, castSuccessful := (definitionURL).(string)
		if !castSuccessful {
			return nil, fmt.Errorf("could not parse rune definition url: %s", definitionURL)
		}
		definitionYaml, yamlErr := yamlSource(definitionURLString)
		if yamlErr != nil {
			return nil, fmt.Errorf("could not parse rune definition at %s", definitionURL)
		}
		definition, mapErr := definitionYaml.Map()
		if mapErr != nil {
			return nil, mapErr
		}
		for k, v := range definition {
			definitions[k] = v
		}
	}

	definitionsYamlBytes, marshalErr := yaml.Marshal(definitions)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return simpleyaml.NewYaml(definitionsYamlBytes)
}

func getRunesRuneDefinitions(runeYaml *simpleyaml.Yaml, stage string) (*simpleyaml.Yaml, error) {
	stageRunes := runeYaml.Get("stages").Get(stage)

	if numRunes, _ := stageRunes.GetArraySize(); numRunes == 0 {
		return nil, fmt.Errorf(
			"Could not find any Runes for stage %s",
			stage,
		)
	}
	return stageRunes, nil
}

func yamlSource(source string) (*simpleyaml.Yaml, error) {
	fileBody, readFileErr := util.ReadFileFromPathOrURL(source)
	if readFileErr != nil {
		return nil, readFileErr
	}

	return simpleyaml.NewYaml(fileBody)
}
