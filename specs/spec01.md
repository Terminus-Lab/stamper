# stamper — Spec

**Status:** Proposal
**Date:** 2026-03-20

---

## Overview

**Stamper** is a standalone CLI tool for human annotation of AI conversation datasets.
Annotators label conversations as `pass / review / fail` with a single keypress — no JSON editing, no setup.

It is a **separate product** from themis. It consumes and produces JSONL, making it compatible with any evaluation system, pandas, or any database.

---

## Command

```bash
stamper -i sampled.jsonl
stamper -i sampled.jsonl -o my_annotations.jsonl
```

**Flags:**

| Flag | Default | Description |
|---|---|---|
| `-i / --input` | required | JSONL file of conversations to annotate |
| `-o / --output` | `{input}_annotated.jsonl` | Annotation output file |

**Output filename rule:** strip `.jsonl` suffix, append `_annotated.jsonl`.
`sampled.jsonl` → `sampled_annotated.jsonl`

---

## Input Schema

JSONL only. Each line must be a valid JSON object with:

| Field | Required | Description |
|---|---|---|
| `conversation_id` | yes | Unique identifier for the conversation |
| `turns` | yes | Array of conversation turns |
| `turns[].query` | yes | The user's message in that turn |
| `turns[].answer` | yes | The agent's response in that turn |
| *(anything else)* | no | Preserved as-is in output (passthrough) |

A turn can contain additional fields — they are ignored but not dropped.
Context fields are not required and not displayed in v1.

---

## Terminal UX

For each conversation, render all turns in full:

```
─────────────────────────────────────────
Conversation 3 / 47  ·  conv-abc123
─────────────────────────────────────────
Turn 1
  User:  What is Python?
  Agent: Python is a high-level programming language...

Turn 2
  User:  Is it hard to learn?
  Agent: Not at all. It has clean syntax...

─────────────────────────────────────────
[p] pass   [r] review   [f] fail   [s] summarize   [x] skip
```

- Single keypress — no Enter needed
- Record written immediately after keypress (safe on Ctrl+C)
- `[x] skip` — conversation is **not written** to output at all
- `[s] summarize` — triggers an LLM call (see Phase 2), appends a summary line inline:

```
─────────────────────────────────────────
Summary: User asked about Python basics. Agent gave accurate, helpful answers.
─────────────────────────────────────────
[p] pass   [r] review   [f] fail   [s] summarize   [x] skip
```

---

## Resume Logic

When the output file already exists, stamper loads all `conversation_id` values from it and silently skips those conversations. Annotators can interrupt (Ctrl+C) and resume later without losing work.

**Choosing a different output file starts fresh** — no conversations are skipped.

---

## Output Format

Original JSON line with `human_annotation` field appended:

```json
{
  "conversation_id": "conv-abc123",
  "turns": [
    {"query": "What is Python?", "answer": "A high-level programming language..."},
    {"query": "Is it hard to learn?", "answer": "Not at all..."}
  ],
  "human_annotation": "pass"
}
```

Skipped conversations produce no output line.

---

## Key Constraints

- **JSONL only** — input and output. Append-safe, human-readable, pandas/DB compatible.
- **Append-safe** — each annotation is flushed immediately. Interrupting mid-session preserves all completed annotations.
- **File-only** — no DB reads or writes in v1.
- **No LLM calls in v1** — `[s] summarize` is a Phase 2 feature.
- **Schema-tolerant** — extra fields on conversations and turns are preserved without modification.

---

## Phases

### Phase 1 — Core annotation CLI
- JSONL input/output
- Single-keypress annotation (p / r / f / x)
- Resume logic via output file deduplication
- All turns displayed in full

### Phase 2 — On-demand summarization
- `[s]` triggers one LLM call for the current conversation
- Summary capped at X tokens, displayed inline above the keypress prompt
- Helps annotators make faster decisions on long conversations
- LLM provider and token cap configurable

### Phase 3 — TUI (stretch)
- `bubbletea`-based terminal UI
- Scrollable turns for long conversations
- Progress bar and end-of-session summary screen

---

## Workflow (themis integration)

```bash

# 1. Annotate (interactive)
stamper -i sampled.jsonl
# writes: sampled_annotated.jsonl

```
