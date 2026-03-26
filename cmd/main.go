package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Terminus-Lab/stamper/internal/annotator"
	"github.com/Terminus-Lab/stamper/internal/config"
	"github.com/Terminus-Lab/stamper/internal/display"
	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/executor"
	"github.com/Terminus-Lab/stamper/internal/logger"
	"github.com/Terminus-Lab/stamper/internal/reader"
	"github.com/Terminus-Lab/stamper/internal/resume"
	"github.com/Terminus-Lab/stamper/internal/tui"
	"github.com/Terminus-Lab/stamper/internal/wire"
	"github.com/Terminus-Lab/stamper/internal/writer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	log := logger.New("info")

	var inputFile string
	var outputFile string

	root := &cobra.Command{
		Use:          "stamper",
		Short:        "Human annotation tool CLI for AI conversation datasets",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputFile == "" {
				return cmd.Help()
			}
			if outputFile == "" {
				outputFile = strings.TrimSuffix(inputFile, ".jsonl") + "_annotated.jsonl"
			}

			log.Info().
				Str("input", inputFile).
				Str("output", outputFile).
				Msg("starting annotation")

			return runAnnotate(inputFile, outputFile, &log)
		},
	}

	root.Flags().StringVarP(&inputFile, "input", "i", "", "JSONL file to annotate (required)")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: {input}_annotated.jsonl)")

	if err := root.Execute(); err != nil {
		log.Fatal().Msg("failed to run stamper")
		os.Exit(1)
	}
}

// tuiEnabled returns false only when STAMPER_TUI is explicitly set to 0, false, or off.
func tuiEnabled() bool {
	switch strings.ToLower(os.Getenv("STAMPER_TUI")) {
	case "0", "false", "off":
		return false
	}
	return true
}

func runAnnotate(inputFile, outputFile string, logger *zerolog.Logger) (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.LoadConfig()

	var exec *executor.Executor
	if cfg.SummarizeEnabled {
		llmClient, err := wire.GetLLMClient(ctx, cfg)
		if err != nil {
			return err
		}
		exec = executor.New(llmClient, cfg, logger)
	}

	res := resume.NewResume(logger)
	rd := reader.NewReader(logger)

	seen, err := res.Load(outputFile)
	if err != nil {
		return err
	}

	all, err := rd.Load(inputFile)
	if err != nil {
		return err
	}

	var remaining []domain.Conversation
	for _, c := range all {
		if !seen[c.ConversationID] {
			remaining = append(remaining, c)
		}
	}

	if len(remaining) == 0 {
		logger.Info().Msg("nothing to annotate")
		return nil
	}

	w, err := writer.New(outputFile)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := w.Close(); cerr != nil {
			logger.Error().Err(cerr).Msg("unable to close the output file")
			err = cerr
		}
	}()

	if tuiEnabled() {
		return runTUI(ctx, remaining, exec, w, cfg.SummarizeEnabled)
	}
	return runPlain(ctx, remaining, exec, w, cfg.SummarizeEnabled, logger)
}

func runTUI(ctx context.Context, remaining []domain.Conversation, exec *executor.Executor, w *writer.Writer, summarizeEnabled bool) error {
	m := tui.New(ctx, remaining, exec, w, summarizeEnabled)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if tm, ok := finalModel.(tui.Model); ok {
		return tm.Err()
	}
	return nil
}

func runPlain(ctx context.Context, remaining []domain.Conversation, exec *executor.Executor, w *writer.Writer, summarizeEnabled bool, logger *zerolog.Logger) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		logger.Info().Msg("interrupted — progress saved")
		os.Exit(0)
	}()

	total := len(remaining)
	for i, conv := range remaining {
		display.Reader(os.Stdout, conv, i+1, total)
		for {
			outcome, err := annotator.ReadKey()
			if err != nil {
				return err
			}
			switch outcome {
			case annotator.OutcomeQuit:
				logger.Info().Msg("interrupted — progress saved")
				return nil
			case annotator.OutcomeSummarize:
				if !summarizeEnabled {
					continue
				}
				fmt.Fprint(os.Stdout, "\n  Summarizing...")
				summary, err := exec.Run(ctx, conv)
				if err != nil {
					fmt.Fprintf(os.Stdout, "\n  Error: %v\n", err)
				} else {
					fmt.Fprintf(os.Stdout, "\r─────────────────────────────────────────\n  Summary: %s\n─────────────────────────────────────────\n", summary)
				}
				continue
			case annotator.OutcomeSkip:
				// no output, move to next conversation
			default:
				if err := w.Append(conv, string(outcome)); err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}
