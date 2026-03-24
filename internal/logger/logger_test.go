package logger

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNew_ValidLevel(t *testing.T) {
	log := New("debug")
	if log.GetLevel() != zerolog.DebugLevel {
		t.Errorf("got %v, want debug", log.GetLevel())
	}
}

func TestNew_InvalidLevelFallsBackToInfo(t *testing.T) {
	log := New("notlevel")
	if log.GetLevel() != zerolog.InfoLevel {
		t.Errorf("got %v, want info", log.GetLevel())
	}
}
