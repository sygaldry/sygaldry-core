package validator

import (
	"testing"
)

// TestRuneRun tests rune.Run
func TestGetValidStage(t *testing.T) {
	stage := "build"
	runesYamlPath := "test_resources/test-runes.yaml"
	_, getValidStageErr := GetValidStage(runesYamlPath, stage)
	if getValidStageErr != nil {
		t.Error("Tried to get stage from runes-test.yaml, got", getValidStageErr)
	}
}
