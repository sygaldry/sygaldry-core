package stage

import (
	"github.com/sygaldry/sygaldry-core/runes"
	"github.com/sygaldry/sygaldry-core/validator"
)

// GetStage will validate and generate a list
// of runes a user needs for a targeted stage
func GetStage(source string, stage string) ([]runes.Rune, error) {
	return validator.GetValidStage(source, stage)
}
