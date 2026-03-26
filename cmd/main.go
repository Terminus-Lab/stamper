package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Terminus-Lab/stamper/internal/annotator"
	"github.com/Terminus-Lab/stamper/internal/display"
	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/logger"
	"github.com/Terminus-Lab/stamper/internal/reader"
	"github.com/Terminus-Lab/stamper/internal/resume"
	"github.com/Terminus-Lab/stamper/internal/writer"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	log := logger.New("info")

	var inputFile string
	var outputFile string

	root := &cobra.Command{
		Use:   "stamper",
		Short: "Human annotation tool CLI for AI conversation datasets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputFile == "" {
				return cmd.Help()
			}
			// Create the output file is not passed by the user
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

	// binds string flag to a variable
	root.Flags().StringVarP(&inputFile, "input", "i", "", "JSONL file to annotate (required)")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: {input}_annotated.jsonl)")

	if err := root.Execute(); err != nil {
		log.Fatal().Msg("failed to run stamper")
		os.Exit(1)
	}
}

func runAnnotate(inputFile, outputFile string, logger *zerolog.Logger) (err error) {
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

	// Handle Ctrl+C outside raw mode (e.g. while display is rendering).
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
		outcome, err := annotator.ReadKey()
		if err != nil {
			return err
		}

		if outcome == annotator.OutcomeQuit {
			logger.Info().Msg("interrupted — progress saved")
			return nil
		}
		if outcome == annotator.OutcomeSkip {
			continue
		}
		if err := w.Append(conv, string(outcome)); err != nil {
			return err
		}
	}

	return nil
}
