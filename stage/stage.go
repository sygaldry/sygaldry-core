package stage

import (
	"fmt"

	"github.com/sygaldry/sygaldry-core/runes"
	"github.com/sygaldry/sygaldry-core/validator"
)

// Stage is a collection of runes
type Stage struct {
	Runes []runes.Rune
}

// GetStage will validate and generate a list
// of runes a user needs for a targeted stage
func GetStage(source string, stage string) (Stage, error) {
	stageRunes, getValidStageErr := validator.GetValidStage(source, stage)
	if getValidStageErr != nil {
		return Stage{}, fmt.Errorf(
			"Could not validate stage %s in runes yaml at %s\n%s",
			stage,
			source,
			getValidStageErr,
		)
	}
	return NewStage(stageRunes), nil
}

// NewStage creates a new stage object
func NewStage(runes []runes.Rune) Stage {
	return Stage{
		Runes: runes,
	}
}

// Run executes all runes in a stage
func (stage *Stage) Run() error {
	for _, stageRune := range stage.Runes {
		stageRuneRunErr := stageRune.Run()
		if stageRuneRunErr != nil {
			return fmt.Errorf(
				"Could not complete Stage execution\n%v",
				stageRuneRunErr,
			)
		}
	}
	return nil
}
