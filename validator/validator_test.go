package validator

import (
	"testing"
)

// TestRuneRun tests rune.Run
func TestYamlSource(t *testing.T) {
	runesYamlPath := "test_resources/test-runes.yaml"
	_, yamlSourceErr := yamlSource(runesYamlPath)
	if yamlSourceErr != nil {
		t.Error("tried to load test-runes.yaml, got", yamlSourceErr)
	}
}

func TestGetRunesConfigsForStage(t *testing.T) {
	runesYamlPath := "test_resources/test-runes.yaml"
	stage := "build"
	runesYaml, yamlSourceErr := yamlSource(runesYamlPath)
	if yamlSourceErr != nil {
		t.Error("tried to load test-runes.yaml, got", yamlSourceErr)
	}
	runesConfigs, getRunesConfigsForStageErr := getRunesConfigsForStage(runesYaml, stage)
	if getRunesConfigsForStageErr != nil {
		t.Errorf("tried to get runeConfigs for stage %s, got %s", stage, getRunesConfigsForStageErr)
	}
	runeConfigsArray, getYamlArrayErr := runesConfigs.Array()
	if getYamlArrayErr != nil {
		t.Errorf("tried to read runeConfigs for stage %s, got %s", stage, getYamlArrayErr)
	}
	if len(runeConfigsArray) != 3 {
		t.Errorf("tried to read all three rune configs for stage %s, found %d", stage, len(runeConfigsArray))
	}
}

func TestGetRuneDefinitions(t *testing.T) {
	runesYamlPath := "test_resources/test-runes.yaml"
	runesYaml, yamlSourceErr := yamlSource(runesYamlPath)
	if yamlSourceErr != nil {
		t.Error("tried to load test-runes.yaml, got", yamlSourceErr)
	}
	runeDefinitions, getRuneDefinitionsErr := getRuneDefinitions(runesYaml)
	if getRuneDefinitionsErr != nil {
		t.Error("tried to get runeDefinitions, got", getRuneDefinitionsErr)
	}
	_, getMapErr := runeDefinitions.Map()
	if getMapErr != nil {
		t.Error("tried to read runeDefinitions as map, got", getMapErr)
	}
	mavenRune := runeDefinitions.Get("MavenRune")
	mavenRuneParams := mavenRune.Get("params")
	mavenRuneParamsList, getArrayErr := mavenRuneParams.Array()
	if getArrayErr != nil {
		t.Error("tried to get MavenRune params, got", getArrayErr)
	}
	if len(mavenRuneParamsList) < 1 {
		t.Errorf("tried to get MavenRune params, got %d params", len(mavenRuneParamsList))
	}
}

func TestValidateRuneConfig(t *testing.T) {
	runesYamlPath := "test_resources/test-runes.yaml"
	stage := "build"
	runesYaml, yamlSourceErr := yamlSource(runesYamlPath)
	if yamlSourceErr != nil {
		t.Error("tried to load test-runes.yaml, got", yamlSourceErr)
	}
	runesConfigs, getRunesConfigsForStageErr := getRunesConfigsForStage(runesYaml, stage)
	if getRunesConfigsForStageErr != nil {
		t.Errorf("tried to get runeConfigs for stage %s, got %s", stage, getRunesConfigsForStageErr)
	}
	runeDefinitions, getRuneDefinitionsErr := getRuneDefinitions(runesYaml)
	if getRuneDefinitionsErr != nil {
		t.Error("tried to get runeDefinitions, got", getRuneDefinitionsErr)
	}
	runeConfig := runesConfigs.GetIndex(0)
	validRuneConfig, validateRuneConfigErr := validateRuneConfig(runeConfig, runeDefinitions)
	if validateRuneConfigErr != nil {
		t.Error("tried to validate runeConfig, got", validateRuneConfigErr)
	}
	if !validRuneConfig {
		t.Error("tried to validate runeConfig, expected true, got", validRuneConfig)
	}

}

func TestGetValidStage(t *testing.T) {
	runesYamlPath := "test_resources/test-runes.yaml"
	stage := "build"
	stageRunes, getValidStageErr := GetValidStage(runesYamlPath, stage)
	if getValidStageErr != nil {
		t.Errorf("tried to read runes for stage %s from %s, got %s", stage, runesYamlPath, getValidStageErr)
	}
	if len(stageRunes) != 3 {
		t.Errorf("tried to read all three rune configs for stage %s, found %d", stage, len(stageRunes))
	}
}
