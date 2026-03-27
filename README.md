# stamper

Human annotation tool for AI conversation datasets. Label conversations as **pass / review / fail** with a single keypress — no JSON editing, no setup.

Reads and writes JSONL, making it compatible with any evaluation pipeline, pandas, or database.

![demo](docs/recording.gif) 

---

## How it works

1. **Load** — reads a JSONL file of conversations (one per line)
2. **Resume** — if an output file already exists, conversations already annotated are silently skipped
3. **Annotate** — renders each conversation turn-by-turn in the terminal and waits for a single keypress
4. **Write** — appends the annotation immediately to the output file (flush on every keypress, safe on Ctrl+C)
5. **Summarize** *(optional)* — press `s` to trigger an LLM call that summarizes the conversation inline before you decide

Each annotation is written as the original JSON line with a `human_annotation` field added. Extra fields on conversations and turns are preserved as-is (passthrough).

---

## Build

**From source:**

```bash
go build -o .bin/stamper ./cmd/
```

> **macOS note:** downloaded binaries may be blocked by Gatekeeper ("not from a trusted developer"). Run once to remove the quarantine flag:
> ```bash
> xattr -d com.apple.quarantine .bin/stamper
> ```

---

## Run

```bash
# Output defaults to sampled_annotated.jsonl
stamper -i sampled.jsonl

# Explicit output file
stamper -i sampled.jsonl -o my_annotations.jsonl

# Without building
go run ./cmd/ -i sampled.jsonl
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-i / --input` | required | JSONL file of conversations to annotate |
| `-o / --output` | `{input}_annotated.jsonl` | Annotation output file |
| `-p / --prompt` | built-in default | Path to a custom prompt template for `[s] summarize` |

---

## Keybindings

| Key | Action |
|---|---|
| `p` | pass |
| `r` | review |
| `f` | fail |
| `s` | summarize via LLM *(requires `STAMPER_SUMMARIZE=true`)* |
| `x` | skip — conversation is not written to output |
| `Ctrl+C` | quit — all completed annotations are preserved |
| `↑ / ↓` | scroll turns *(TUI mode)* |

---

## Resume

Interrupt at any time with `Ctrl+C`. Re-run the same command — stamper reads the output file on startup and silently skips already-annotated conversations.

Choosing a different output file starts fresh.

---

## Configuration

Copy `.env.example` to `.env` and fill in the values. A `.env` file is optional — any variable can also be exported directly in the shell.

```bash
cp .env.example .env
```

| Variable | Default | Description |
|---|---|---|
| `LLM_FAMILY` | `openai_platform` | LLM provider: `openai_platform`, `openai` (Azure), `ollama` |
| `MODEL_ID` | — | Model name, e.g. `gpt-4o-mini`, `llama3` |
| `OPEN_AI_KEY` | — | OpenAI or Azure API key |
| `AZURE_OPENAI_ENDPOINT` | — | Azure OpenAI endpoint URL |
| `OLLAMA_BASE_URL` | `http://localhost:11434/v1` | Ollama base URL |
| `MODEL_MAX_TOKENS` | `1000` | Max tokens for summarize calls |
| `MODEL_TEMPERATURE` | `0.0` | Sampling temperature |
| `STAMPER_SUMMARIZE` | `false` | Enable `[s] summarize` — no LLM client is created when false |
| `STAMPER_TUI` | `true` | Set to `false` for plain terminal mode (no bubbletea) |

### Prompt template

The summarize prompt is built into the binary. To customize it, edit `conf/summarize_prompt.tmpl` (included in every release archive) and pass it at runtime:

```bash
stamper -i sampled.jsonl -p conf/summarize_prompt.tmpl
```

The template has access to `.Turns` (array of `Query` / `Answer`) and the `inc` helper to produce 1-based turn numbers.

---

## Input format

JSONL — one conversation per line:

```json
{"conversation_id": "conv-001", "turns": [{"query": "What is Python?", "answer": "A high-level language..."}]}
{"conversation_id": "conv-002", "turns": [{"query": "Is it hard?", "answer": "Not at all..."}]}
```

| Field | Required | Description |
|---|---|---|
| `conversation_id` | yes | Unique identifier |
| `turns` | yes | Array of turns |
| `turns[].query` | yes | User message |
| `turns[].answer` | yes | Agent response |
| *(anything else)* | no | Preserved as-is in output |

## Output format

Original JSON with `human_annotation` appended:

```json
{"conversation_id": "conv-001", "turns": [...], "human_annotation": "pass"}
```

Skipped conversations (`x`) produce no output line.

---

## Development

```bash
# Run tests
go test ./...

# Lint
golangci-lint run

# Sample dataset (20 conversations)
stamper -i resources/sampled.jsonl
```
