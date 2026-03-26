package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/executor"
	"github.com/Terminus-Lab/stamper/internal/writer"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// fixedLines is the number of lines outside the scrollable viewport:
// progress bar + blank line + separator + header + separator + separator + footer
const fixedLines = 7

const separator = "─────────────────────────────────────────"

type summaryMsg struct {
	text string
	err  error
}

type Model struct {
	conversations []domain.Conversation
	index         int
	total         int
	viewport      viewport.Model
	progress      progress.Model
	writer        *writer.Writer
	exec          *executor.Executor
	ctx           context.Context
	ready         bool
	loading       bool
	summary       string
	err           error
}

func New(ctx context.Context, conversations []domain.Conversation, exec *executor.Executor, w *writer.Writer) Model {
	return Model{
		conversations: conversations,
		total:         len(conversations),
		progress:      progress.New(progress.WithDefaultGradient()),
		writer:        w,
		exec:          exec,
		ctx:           ctx,
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
			m.viewport.SetContent(renderContent(m.conversations[m.index], ""))
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

	case summaryMsg:
		m.loading = false
		if msg.err != nil {
			m.summary = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.summary = msg.text
		}
		m.viewport.SetContent(renderContent(m.conversations[m.index], m.summary))
		m.viewport.GotoBottom()
		return m, nil

	case tea.KeyMsg:
		if !m.ready {
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			if m.loading {
				return m, nil
			}
			m.loading = true
			m.summary = ""
			conv := m.conversations[m.index]
			return m, fetchSummary(m.ctx, m.exec, conv)
		case "p", "r", "f", "x":
			return m.annotate(msg.String())
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func fetchSummary(ctx context.Context, exec *executor.Executor, conv domain.Conversation) tea.Cmd {
	return func() tea.Msg {
		text, err := exec.Run(ctx, conv)
		return summaryMsg{text: text, err: err}
	}
}

func (m Model) annotate(key string) (tea.Model, tea.Cmd) {
	outcome := outcomeFor(key)
	if outcome != "skip" {
		if err := m.writer.Append(m.conversations[m.index], outcome); err != nil {
			m.err = err
			return m, tea.Quit
		}
	}
	m.index++
	if m.index >= m.total {
		return m, tea.Quit
	}
	m.summary = ""
	m.loading = false
	m.viewport.SetContent(renderContent(m.conversations[m.index], ""))
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

	summarizeLabel := "[s] summarize"
	if m.loading {
		summarizeLabel = "[s] summarizing..."
	}
	footer := fmt.Sprintf("[p] pass   [r] review   [f] fail   %s   [x] skip   [↑↓] scroll%s", summarizeLabel, scrollHint)

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

func renderContent(conv domain.Conversation, summary string) string {
	var s strings.Builder
	for i, turn := range conv.Turns {
		fmt.Fprintf(&s, "Turn %d\n", i+1)
		fmt.Fprintf(&s, "  User:  %s\n\n", turn.Query)
		fmt.Fprintf(&s, "  Agent: %s\n\n", turn.Answer)
	}
	if summary != "" {
		s.WriteString(separator + "\n")
		fmt.Fprintf(&s, "  Summary: %s\n", summary)
		s.WriteString(separator + "\n")
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
	case "x":
		return "skip"
	}
	return ""
}
