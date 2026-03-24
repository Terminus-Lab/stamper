package annotator

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type Outcome string

const (
	OutcomePass   Outcome = "pass"
	OutcomeReview Outcome = "review"
	OutcomeSkip   Outcome = "skip"
	OutcomeFail   Outcome = "fail"
)

func ReadKey() (Outcome, error) {
	// Switch stdin to raw mode

	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("stdin is not a terminal (TTY)")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}

	defer term.Restore(fd, oldState)
	return ReadKeyFrom(os.Stdin)
}

func ReadKeyFrom(r io.Reader) (Outcome, error) {
	buf := make([]byte, 1)
	for {
		if _, err := r.Read(buf); err != nil {
			return "", err
		}
		switch buf[0] {
		case 'p':
			return OutcomePass, nil
		case 'r':
			return OutcomeReview, nil
		case 'f':
			return OutcomeFail, nil
		case 'x':
			return OutcomeSkip, nil
		default:
			// wait for a valid key
		}
	}
}
