package resume

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Resume struct {
	Logger *zerolog.Logger
}

func NewResume(logger *zerolog.Logger) *Resume {
	return &Resume{
		Logger: logger,
	}
}

func (r *Resume) Load(path string) (map[string]bool, error) {
	annotated := make(map[string]bool)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return annotated, nil
	}
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil {
			r.Logger.Error().Err(cerr).Msg("unable to close resume file")
		}
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var c struct {
			ConversationID string `json:"conversation_id"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			return nil, fmt.Errorf("invalid line in output file: %w", err)
		}
		annotated[c.ConversationID] = true
	}

	return annotated, scanner.Err()
}
