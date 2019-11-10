package validator

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

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
	runeConfigsLen, getRuneConfigsLenErr := runeConfigs.GetArraySize()

	if getRuneConfigsLenErr != nil {
		return nil, fmt.Errorf("Could not get count of runes for stage %s", stage)
	}

	validStageRunes := make([]runes.Rune, runeConfigsLen)

	for index := 0; index < runeConfigsLen; index++ {
		runeConfig := runeConfigs.GetIndex(index)
		if runeConfig != nil {
			if valid, _ := validateRuneConfig(runeConfig, runeDefinitions); valid {
				validRune, buildValidRuneErr := buildValidRune(runeConfig, runeDefinitions)
				if buildValidRuneErr != nil {
					return nil, errors.New("Could not build rune")
				}
				validStageRunes[index] = validRune
			}
		} else {
			return nil, fmt.Errorf("Could not parse runes for stage %s", stage)
		}
	}

	return validStageRunes, nil
}

func buildValidRune(runeConfig *simpleyaml.Yaml, runeDefinitions *simpleyaml.Yaml) (runes.Rune, error) {
	runeConfigDefinitionName, getStringErr := runeConfig.Get("definition").String()
	if getStringErr != nil {
		return runes.Rune{}, errors.New("Could not read rune definition name")
	}

	runeDefinition := runeDefinitions.Get(runeConfigDefinitionName)

	runeDefinitionImageTemplate, getImageTemplateErr := runeDefinition.Get("values").Get("Image").String()
	if getImageTemplateErr != nil {
		return runes.Rune{}, fmt.Errorf("Could not read Image value for rune definition %s", runeConfigDefinitionName)
	}
	runeDefinitionEnvTemplates, getEnvTemplateErr := util.ConvertYamlListToStringList(runeDefinition.Get("values").Get("Env"))
	if getEnvTemplateErr != nil {
		runeDefinitionEnvTemplates = []string{}
	}
	runeDefinitionVolumesTemplates, getVolumesTemplateErr := util.ConvertYamlListToStringList(runeDefinition.Get("values").Get("Volumes"))
	if getVolumesTemplateErr != nil {
		runeDefinitionVolumesTemplates = []string{}
	}

	runeImage, imageBuildStringForTemplateErr := buildStringsForTemplate(runeDefinitionImageTemplate, runeConfig)
	if imageBuildStringForTemplateErr != nil {
		return runes.Rune{}, imageBuildStringForTemplateErr
	}
	runeEnv, envBuildStringsForTemplatesErr := buildStringsForTemplates(runeDefinitionEnvTemplates, runeConfig)
	if envBuildStringsForTemplatesErr != nil {
		return runes.Rune{}, envBuildStringsForTemplatesErr
	}
	runeVolumes, volumesBuildStringsForTemplatesErr := buildStringsForTemplates(runeDefinitionVolumesTemplates, runeConfig)
	if volumesBuildStringsForTemplatesErr != nil {
		return runes.Rune{}, volumesBuildStringsForTemplatesErr
	}

	return runes.NewRune(runeImage, runeEnv, runeVolumes)
}

func buildStringsForTemplate(runeValueTemplate string, runeConfig *simpleyaml.Yaml) (string, error) {
	runeConfigMap, getRuneConfigMapErr := runeConfig.Map()
	if getRuneConfigMapErr != nil {
		return "", errors.New("Could not make map out of runeConfig")
	}
	valueTemplate, newValueTemplateErr := template.New("value").Parse(runeValueTemplate)
	if newValueTemplateErr != nil {
		return "", fmt.Errorf("Cound not create go template for: %s", runeValueTemplate)
	}
	var valueBytes bytes.Buffer
	valueTemplate.Execute(&valueBytes, runeConfigMap)
	return valueBytes.String(), nil
}

func buildStringsForTemplates(runeValueTemplates []string, runeConfig *simpleyaml.Yaml) ([]string, error) {
	runeConfigMap, getRuneConfigMapErr := runeConfig.Map()
	if getRuneConfigMapErr != nil {
		return nil, errors.New("Could not make map out of runeConfig")
	}
	numValues := len(runeValueTemplates)
	valuesBytes := make([]bytes.Buffer, numValues)
	valuesStrings := make([]string, numValues)
	for index, runeValueTemplate := range runeValueTemplates {
		valueTemplate, newTemplateErr := template.New(fmt.Sprintf("value-%d", index)).Parse(runeValueTemplate)
		if newTemplateErr != nil {
			return nil, fmt.Errorf("Cound not create go template for: %s", runeValueTemplate)
		}
		valueTemplate.Execute(&valuesBytes[index], runeConfigMap)
		valuesStrings[index] = valuesBytes[index].String()
	}
	return valuesStrings, nil
}

func validateRuneConfig(runeConfig *simpleyaml.Yaml, runeDefinitions *simpleyaml.Yaml) (bool, error) {
	runeConfigDefinitionName, getStringErr := runeConfig.Get("definition").String()
	if getStringErr != nil {
		return false, errors.New("Could not read rune definition name")
	}
	runeDefinition := runeDefinitions.Get(runeConfigDefinitionName)
	runeConfigParams, getMapKeysErr := runeConfig.GetMapKeys()
	if getMapKeysErr != nil {
		return false, fmt.Errorf("could not read keys for %s rune", runeConfigDefinitionName)
	}

	return validateRuneConfigParams(runeConfigParams, runeDefinition)
}

func validateRuneConfigParams(runeConfigParams []string, runeDefinition *simpleyaml.Yaml) (bool, error) {
	runeDefinitionParams := runeDefinition.Get("params")
	runeDefinitionParamsList, getStringErr := util.ConvertYamlListToStringList(runeDefinitionParams)
	if getStringErr != nil {
		errors.New("could not convert rune definition params to list of strings")
	}
	return util.ListContainsListStrings(runeConfigParams, runeDefinitionParamsList), nil
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

func yamlSource(source string) (*simpleyaml.Yaml, error) {
	fileBody, readFileErr := util.ReadFileFromPathOrURL(source)
	if readFileErr != nil {
		return nil, readFileErr
	}

	return simpleyaml.NewYaml(fileBody)
}
