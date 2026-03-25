package writer

import (
	"os"
	"testing"
)

func TestAppend(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "out*.jsonl")

	f.Close()

}
