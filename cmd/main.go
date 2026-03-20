package main

import (
	"os"
	"strings"

	"github.com/Terminus-Lab/stamper/internal/logger"
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
			// Create the output file is not passed by the user
			if outputFile == "" {
				outputFile = strings.TrimSuffix(inputFile, ".jsonl") + "_annotated.jsonl"
			}

			log.Info().
				Str("input", inputFile).
				Str("output", outputFile).
				Msg("Starting annotation")

			return nil
		},
	}

	// binds string flag to a variable
	root.Flags().StringVarP(&inputFile, "input", "i", "", "JSONL file to annotate (required)")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: {input}_annotated.jsonl)")
	root.MarkFlagRequired("input") //makes Cobra emit an error automatically if -i is missing

	if err := root.Execute(); err != nil {
		log.Fatal().Msg("Failed to run stamper")
		os.Exit(1)
	}
}
