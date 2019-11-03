package runes_test

import (
	"testing"

	"github.com/sygaldry/sygaldry-core/runes"
)

// TestRuneRun tests rune.Run
func TestRuneRun(t *testing.T) {
	myFirstRune, error := runes.NewRune("docker.io/library/hello-world", []string{"TEST=test"}, []string{})
	if error != nil {
		t.Error("Expected to create a Rune, got", error)
	}
	if error := myFirstRune.Run(); error != nil {
		t.Error("Expected success, got ", error)
	}
}
