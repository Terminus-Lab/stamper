package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Terminus-Lab/stamper/internal/domain"
)

type Writer struct {
	F *os.File
}

func New(path string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open output file: %w", err)
	}

	return &Writer{F: f}, nil
}

func (w *Writer) Append(conv domain.Conversation, annotation string) error {
	conv.Annotation = annotation
	line, err := json.Marshal(&conv)
	if err != nil {
		return fmt.Errorf("marshal conversation: %w", err)
	}

	if _, err := fmt.Fprintf(w.F, "%s\n", line); err != nil {
		return fmt.Errorf("write line: %w", err)
	}

	return w.F.Sync()
}

func (w *Writer) Close() error {
	return w.F.Close()
}
