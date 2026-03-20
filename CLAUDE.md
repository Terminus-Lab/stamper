# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

**Stamper** is a Go CLI tool for human annotation of AI conversation datasets. Annotators label conversations as `pass / review / fail` with a single keypress. It reads and writes JSONL, making it compatible with any evaluation pipeline.

The full product spec is in `specs/spec01.md` — read it before making significant changes.

## Commands

```bash
# Build
go build ./cmd/...

# Run
go run ./cmd/main.go -i sampled.jsonl
go run ./cmd/main.go -i sampled.jsonl -o my_annotations.jsonl

# Test
go test ./...

# Run a single test
go test ./... -run TestFunctionName

# Lint (install golangci-lint if not present)
golangci-lint run
```

## Architecture

The project uses [Cobra](https://github.com/spf13/cobra) for CLI argument parsing. Currently `cmd/main.go` is a stub — the intended structure from the spec is:

- **Input parsing**: Read JSONL line-by-line; each line is a conversation with `conversation_id` and `turns[]` (each turn has `query` and `answer`). Extra fields must be preserved passthrough.
- **Resume logic**: On startup, if the output file already exists, load all `conversation_id` values from it and skip those conversations silently.
- **Interactive loop**: For each conversation, render all turns to the terminal, then wait for a single keypress (`p` / `r` / `f` / `x`). No Enter needed.
- **Output writing**: Append each annotation immediately after keypress (flush to disk). Skipped conversations (`x`) produce no output line. Output is the original JSON with `"human_annotation"` field added.
- **Output filename default**: strip `.jsonl` suffix, append `_annotated.jsonl` (e.g., `sampled.jsonl` → `sampled_annotated.jsonl`).
