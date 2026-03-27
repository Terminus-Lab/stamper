package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/Terminus-Lab/stamper/internal/domain"
)

const separator = "─────────────────────────────────────────\n"

func Reader(w io.Writer, conv domain.Conversation, index, total int) {
	var s strings.Builder

	s.WriteString(separator)
	fmt.Fprintf(&s, "Conversation %d / %d  - %d turns\n", index, total, len(conv.Turns))
	s.WriteString(separator)
	for i, turn := range conv.Turns {
		fmt.Fprintf(&s, "Turn %d\n", i+1)
		fmt.Fprintf(&s, "\t User: %s\n", turn.UserQuery)
		fmt.Fprintf(&s, "\t Answer: %s\n", turn.Answer)
	}
	s.WriteString(separator)
	s.WriteString("[p] pass   [r] review   [f] fail   [s]summarize   [x] skip\n")
	_, _ = w.Write([]byte(s.String()))
}
