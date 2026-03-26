# Changelog

## [Unreleased]

### Phase 2 — On-demand summarization

- **`[s] summarize`** — press `s` on any conversation to trigger an LLM call that generates a 2-3 sentence summary inline, helping annotators make faster decisions on long conversations
- **Multi-provider LLM support** — OpenAI Platform, Azure OpenAI, and Ollama via a unified client interface
- **`STAMPER_SUMMARIZE` flag** — summarization is opt-in (`false` by default); no LLM client is created when disabled
- **`.env` file support** — configuration via `.env` file or exported shell variables; `.env` is optional

---

## [0.1.0] — Phase 1 — Core annotation CLI

### Features

- **Single-keypress annotation** — label conversations `pass / review / fail / skip` with no Enter key required
- **JSONL input/output** — one conversation per line; compatible with pandas, any database, or evaluation pipeline
- **Resume logic** — on startup, already-annotated conversations are silently skipped based on the output file; interrupt and resume at any time safely
- **Append-safe writes** — each annotation is flushed to disk immediately after keypress; Ctrl+C preserves all completed work
- **Schema passthrough** — extra fields on conversations and turns are preserved as-is in the output
- **TUI mode** — scrollable bubbletea interface with a progress bar (disable with `STAMPER_TUI=false` for plain terminal)
- **Output filename default** — `sampled.jsonl` → `sampled_annotated.jsonl`
