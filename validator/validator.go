package validator

import (
	"bytes"
	"fmt"
	"strings"
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
		return nil, fmt.Errorf("Could not parse runes yaml at: %s\n%v", source, newYamlErr)
	}

	runeDefinitions, getRuneDefinitionsErr := getRuneDefinitions(runeYaml)
	if getRuneDefinitionsErr != nil {
		return nil, fmt.Errorf("Could not read rune definitions in runes yaml at: %s\n%v", source, getRuneDefinitionsErr)
	}
	runeConfigs, getRuneConfigsErr := getRunesConfigsForStage(runeYaml, stage)
	if getRuneConfigsErr != nil {
		return nil, fmt.Errorf("Could not read rune configs for stage %s in runes yaml at %s\n%v", stage, source, getRuneConfigsErr)
	}
	runeConfigsLen, getRuneConfigsLenErr := runeConfigs.GetArraySize()

	if getRuneConfigsLenErr != nil {
		return nil, fmt.Errorf("Could not get count of runes for stage %s\n%v", stage, getRuneConfigsLenErr)
	}

	validStageRunes := make([]runes.Rune, runeConfigsLen)

	for index := 0; index < runeConfigsLen; index++ {
		runeConfig := runeConfigs.GetIndex(index)
		if runeConfig != nil {
			runeConfigValid, validateRuneConfigErr := validateRuneConfig(runeConfig, runeDefinitions)
			if validateRuneConfigErr != nil {
				return nil, fmt.Errorf(
					"Could not validate rune config at index %d of stage %s in runes yaml at %s\n%v",
					index,
					stage,
					source,
					validateRuneConfigErr,
				)
			}
			if runeConfigValid {
				validRune, buildValidRuneErr := buildValidRune(runeConfig, runeDefinitions)
				if buildValidRuneErr != nil {
					return nil, fmt.Errorf(
						"Could not build rune at index %d of stage %s in runes yaml at %s\n%v",
						index,
						stage,
						source,
						buildValidRuneErr,
					)
				}
				validStageRunes[index] = validRune
			}
		} else {
			return nil, fmt.Errorf(
				"Could not parse runes for stage %s in runes yaml at %s",
				stage,
				source,
			)
		}
	}

	return validStageRunes, nil
}

func buildValidRune(runeConfig *simpleyaml.Yaml, runeDefinitions *simpleyaml.Yaml) (runes.Rune, error) {
	runeConfigDefinitionName, getStringErr := runeConfig.Get("definition").String()
	if getStringErr != nil {
		return runes.Rune{}, fmt.Errorf("Could not parse definition field of rune\n%v", getStringErr)
	}

	runeDefinition := runeDefinitions.Get(runeConfigDefinitionName)

	runeDefinitionImageTemplate, getImageTemplateErr := runeDefinition.Get("values").Get("Image").String()
	if getImageTemplateErr != nil {
		return runes.Rune{}, fmt.Errorf(
			"Could not read Image value of %s rune definition\n%v",
			runeConfigDefinitionName,
			getImageTemplateErr,
		)
	}
	runeDefinitionEnvTemplates, getEnvTemplateErr := util.ConvertYamlListToStringList(runeDefinition.Get("values").Get("Env"))
	if getEnvTemplateErr != nil {
		fmt.Printf(
			"WARNING: Could not read Env values of %s rune definition (this may be expected)\n%v",
			runeConfigDefinitionName,
			getEnvTemplateErr,
		)
		runeDefinitionEnvTemplates = []string{}
	}
	runeDefinitionVolumesTemplates, getVolumesTemplateErr := util.ConvertYamlListToStringList(runeDefinition.Get("values").Get("Volumes"))
	if getVolumesTemplateErr != nil {
		fmt.Printf(
			"WARNING: Could not read Volumes values of %s rune definition (this may be expected)\n%v",
			runeConfigDefinitionName,
			getVolumesTemplateErr,
		)
		runeDefinitionVolumesTemplates = []string{}
	}

	runeImage, imageBuildStringForTemplateErr := buildStringsForTemplate(runeDefinitionImageTemplate, runeConfig)
	if imageBuildStringForTemplateErr != nil {
		return runes.Rune{}, fmt.Errorf(
			"Could not evalute template for %s definition from Image value\n%v",
			runeConfigDefinitionName,
			imageBuildStringForTemplateErr,
		)
	}
	runeEnv, envBuildStringsForTemplatesErr := buildStringsForTemplates(runeDefinitionEnvTemplates, runeConfig)
	if envBuildStringsForTemplatesErr != nil {
		return runes.Rune{}, fmt.Errorf(
			"Could not evalute template for %s definition from Env value(s)\n%v",
			runeConfigDefinitionName,
			envBuildStringsForTemplatesErr,
		)
	}
	runeVolumes, volumesBuildStringsForTemplatesErr := buildStringsForTemplates(runeDefinitionVolumesTemplates, runeConfig)
	if volumesBuildStringsForTemplatesErr != nil {
		return runes.Rune{}, fmt.Errorf(
			"Could not evalute template for %s definition from Volumes value(s)\n%v",
			runeConfigDefinitionName,
			volumesBuildStringsForTemplatesErr,
		)
	}

	return runes.NewRune(runeImage, runeEnv, runeVolumes)
}

func buildStringsForTemplate(runeValueTemplate string, runeConfig *simpleyaml.Yaml) (string, error) {
	runeConfigMap, getRuneConfigMapErr := runeConfig.Map()
	if getRuneConfigMapErr != nil {
		return "", fmt.Errorf(
			"Could not convert rune config to Map while parsing %s\n%v",
			runeValueTemplate,
			getRuneConfigMapErr,
		)
	}
	valueTemplate, newValueTemplateErr := template.New("value").Parse(runeValueTemplate)
	if newValueTemplateErr != nil {
		return "", fmt.Errorf(
			"Cound not create go template while parsing %s\n%v",
			runeValueTemplate,
			newValueTemplateErr,
		)
	}
	var valueBytes bytes.Buffer
	valueTemplate.Execute(&valueBytes, runeConfigMap)
	return valueBytes.String(), nil
}

func buildStringsForTemplates(runeValueTemplates []string, runeConfig *simpleyaml.Yaml) ([]string, error) {
	runeConfigMap, getRuneConfigMapErr := runeConfig.Map()
	if getRuneConfigMapErr != nil {
		return nil, fmt.Errorf(
			"Could not convert rune config to Map while parsing %s\n%v",
			strings.Join(runeValueTemplates[:], ", "),
			getRuneConfigMapErr,
		)
	}
	numValues := len(runeValueTemplates)
	valuesBytes := make([]bytes.Buffer, numValues)
	valuesStrings := make([]string, numValues)
	for index, runeValueTemplate := range runeValueTemplates {
		valueTemplate, newTemplateErr := template.New(fmt.Sprintf("value-%d", index)).Parse(runeValueTemplate)
		if newTemplateErr != nil {
			return nil, fmt.Errorf(
				"Cound not create go template while parsing %s\n%v",
				runeValueTemplate,
				newTemplateErr,
			)
		}
		valueTemplate.Execute(&valuesBytes[index], runeConfigMap)
		valuesStrings[index] = valuesBytes[index].String()
	}
	return valuesStrings, nil
}

func validateRuneConfig(runeConfig *simpleyaml.Yaml, runeDefinitions *simpleyaml.Yaml) (bool, error) {
	runeConfigDefinitionName, getStringErr := runeConfig.Get("definition").String()
	if getStringErr != nil {
		return false, fmt.Errorf("Could not parse definition field of rune\n%v", getStringErr)
	}
	runeDefinition := runeDefinitions.Get(runeConfigDefinitionName)
	runeConfigParams, getMapKeysErr := runeConfig.GetMapKeys()
	if getMapKeysErr != nil {
		return false, fmt.Errorf(
			"Could not read keys for %s rune\n%v",
			runeConfigDefinitionName,
			getMapKeysErr,
		)
	}

	return validateRuneConfigParams(runeConfigParams, runeDefinition)
}

func validateRuneConfigParams(runeConfigParams []string, runeDefinition *simpleyaml.Yaml) (bool, error) {
	runeDefinitionParams := runeDefinition.Get("params")
	runeDefinitionParamsList, getStringErr := util.ConvertYamlListToStringList(runeDefinitionParams)
	if getStringErr != nil {
		return false, fmt.Errorf(
			"Could not convert rune definition params to list of strings\n%v",
			getStringErr,
		)
	}
	return util.ListContainsListStrings(runeConfigParams, runeDefinitionParamsList), nil
}

func getRunesConfigsForStage(runeYaml *simpleyaml.Yaml, stage string) (*simpleyaml.Yaml, error) {
	stageRunes := runeYaml.Get("stages").Get(stage)
	numRunes, getArraySizeError := stageRunes.GetArraySize()
	if getArraySizeError != nil {
		return nil, fmt.Errorf(
			"Could not find any Runes for stage %s\n%v",
			stage,
			getArraySizeError,
		)
	}
	if numRunes == 0 {
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
		return nil, fmt.Errorf(
			"Could not find rune definitions\n%v",
			getDefinitionURLsErr,
		)
	}

	definitions := map[interface{}]interface{}{}

	for _, definitionURL := range definitionURLsArray {
		definitionURLString, castSuccessful := (definitionURL).(string)
		if !castSuccessful {
			return nil, fmt.Errorf("Could not parse rune definition at: %s", definitionURL)
		}
		definitionYaml, yamlSourceErr := yamlSource(definitionURLString)
		if yamlSourceErr != nil {
			return nil, fmt.Errorf(
				"Could not parse rune definition at: %s\n%v",
				definitionURL,
				yamlSourceErr,
			)
		}
		definition, mapErr := definitionYaml.Map()
		if mapErr != nil {
			return nil, fmt.Errorf(
				"Could not build map from rune definition at %s\n%v",
				definitionURL,
				mapErr,
			)
		}
		for k, v := range definition {
			definitions[k] = v
		}
	}

	definitionsYamlBytes, marshalErr := yaml.Marshal(definitions)
	if marshalErr != nil {
		return nil, fmt.Errorf(
			"Could not make yaml from collection of rune definitions\n%v",
			marshalErr,
		)
	}

	return simpleyaml.NewYaml(definitionsYamlBytes)
}

func yamlSource(source string) (*simpleyaml.Yaml, error) {
	fileBody, readFileErr := util.ReadFileFromPathOrURL(source)
	if readFileErr != nil {
		return nil, fmt.Errorf(
			"Could not read yaml at %s\n%v",
			source,
			readFileErr,
		)
	}

	return simpleyaml.NewYaml(fileBody)
}
