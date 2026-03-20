package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/rs/zerolog"
)

type Reader struct {
	Logger *zerolog.Logger
}

func NewReader(logger *zerolog.Logger) *Reader {
	return &Reader{
		Logger: logger,
	}
}

func (r *Reader) Load(path string) ([]domain.Conversation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			r.Logger.Error().Err(err).Msg("Unable to close the file")
		}
	}()

	var conversations []domain.Conversation
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		var c domain.Conversation
		if err := json.Unmarshal(line, &c); err != nil {
			return nil, fmt.Errorf("invalid JSON line: %w", err)
		}
		conversations = append(conversations, c)
	}

	return conversations, scanner.Err()

}
