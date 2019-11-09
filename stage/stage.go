package stage

import (
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
	return NewStage(stageRunes), getValidStageErr
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
			panic(stageRuneRunErr)
		}
	}
	return nil
}
