# stamper

Human annotation tool for AI conversation datasets. Single-keypress labeling, resumable sessions, JSONL in/out.

## Quick start

```bash
# Build
go build -o .bin/stamper ./cmd/main.go

# Run with the sample dataset (20 conversations)
.bin/stamper -i resource/sampled.jsonl

# Or run without building first
go run ./cmd/main.go -i sampled.jsonl
```

Output is written to `sampled_annotated.jsonl` by default.

## Flags

| Flag | Default | Description |
|---|---|---|
| `-i / --input` | required | JSONL file of conversations to annotate |
| `-o / --output` | `{input}_annotated.jsonl` | Annotation output file |

## Keybindings

| Key | Action |
|---|---|
| `p` | pass |
| `r` | review |
| `f` | fail |
| `x` | skip (no output written) |

## Resume

If you interrupt with Ctrl+C, re-run the same command. Already-annotated conversations are skipped automatically based on the output file.

## Input format

JSONL — one conversation per line:

```json
{"conversation_id": "conv-001", "turns": [{"query": "What is Python?", "answer": "A high-level language..."}]}
```

## Output format

Original JSON with `human_annotation` appended:

```json
{"conversation_id": "conv-001", "turns": [...], "human_annotation": "pass"}
```
