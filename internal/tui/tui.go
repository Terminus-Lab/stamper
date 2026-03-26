package tui

import (
	"fmt"
	"strings"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/writer"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// fixedLines is the number of lines outside the scrollable viewport:
// progress bar + blank line + separator + header + separator + separator + footer
const fixedLines = 7

const separator = "─────────────────────────────────────────"

type Model struct {
	conversations []domain.Conversation
	index         int
	total         int
	viewport      viewport.Model
	progress      progress.Model
	writer        *writer.Writer
	ready         bool
	err           error
}

func New(conversations []domain.Conversation, w *writer.Writer) Model {
	return Model{
		conversations: conversations,
		total:         len(conversations),
		progress:      progress.New(progress.WithDefaultGradient()),
		writer:        w,
	}
}

// Err returns any write error encountered during annotation.
func (m Model) Err() error {
	return m.err
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		vpHeight := msg.Height - fixedLines
		if vpHeight < 1 {
			vpHeight = 1
		}
		if !m.ready {
			m.viewport = viewport.New(msg.Width, vpHeight)
			m.viewport.SetContent(renderContent(m.conversations[m.index]))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = vpHeight
		}
		return m, nil

	case progress.FrameMsg:
		pm, cmd := m.progress.Update(msg)
		m.progress = pm.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		if !m.ready {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "p", "r", "f", "s", "x":
			return m.annotate(msg.String())
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) annotate(key string) (tea.Model, tea.Cmd) {
	outcome := outcomeFor(key)
	if outcome != "skip" {
		if err := m.writer.Append(m.conversations[m.index], outcome); err != nil {
			m.err = err
			return m, tea.Quit
		}
	}
	if outcome == "summarize" {
		//update llm call
		if err := m.writer.Append(m.conversations[m.index], outcome); err != nil {
			m.err = err
			return m, tea.Quit
		}
	}
	m.index++
	if m.index >= m.total {
		return m, tea.Quit
	}
	m.viewport.SetContent(renderContent(m.conversations[m.index]))
	m.viewport.GotoTop()
	cmd := m.progress.SetPercent(float64(m.index) / float64(m.total))
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}

	conv := m.conversations[m.index]
	header := fmt.Sprintf("Conversation %d / %d  ·  %s  ·  %d turns",
		m.index+1, m.total, conv.ConversationID, len(conv.Turns))

	scrollHint := ""
	if m.viewport.TotalLineCount() > m.viewport.Height {
		scrollHint = fmt.Sprintf("  %3.f%% ↕", m.viewport.ScrollPercent()*100)
	}
	footer := fmt.Sprintf("[p] pass   [r] review   [f] fail   [s]summarize   [x] skip   [↑↓] scroll%s", scrollHint)

	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s\n%s",
		m.progress.View(),
		separator,
		header,
		separator,
		m.viewport.View(),
		separator,
		footer,
	)
}

func renderContent(conv domain.Conversation) string {
	var s strings.Builder
	for i, turn := range conv.Turns {
		fmt.Fprintf(&s, "Turn %d\n", i+1)
		fmt.Fprintf(&s, "  User:  %s\n\n", turn.Query)
		fmt.Fprintf(&s, "  Agent: %s\n\n", turn.Answer)
	}
	return s.String()
}

func outcomeFor(key string) string {
	switch key {
	case "p":
		return "pass"
	case "r":
		return "review"
	case "f":
		return "fail"
	case "s":
		return "summarize"
	case "x":
		return "skip"
	}
	return ""
}
