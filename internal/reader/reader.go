package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Terminus-Lab/stamper/internal/domain"
)

func Load(path string) ([]domain.Conversation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

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
