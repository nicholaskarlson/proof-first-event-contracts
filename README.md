# proof-first-event-contracts

**Repo A (Book 2): event contracts** — fixture-backed parsing + filtering for Cloud Storage / Eventarc-style CloudEvents.

![ci](https://github.com/nicholaskarlson/proof-first-event-contracts/actions/workflows/ci.yml/badge.svg)
![license](https://img.shields.io/badge/license-MIT-blue.svg)

> **Book:** *Proof-First Pipelines in the Cloud* (Book 2)  
> This repo is **Repo A (Repo 2 of 4)**. The exact code referenced in the manuscript is tagged **[`book2-v1`](https://github.com/nicholaskarlson/proof-first-event-contracts/tree/book2-v1)**.

**Goal:** turn “incoming storage event” into deterministic artifacts you can review as diffs:

- `decision.json` — run vs. ignore (with a stable reason)
- `object_ref.json` — parsed bucket + object name (and URL-unescaped variant)
- `error.txt` — expected-fail contract (stable error text, ends with a newline)

**Go baseline:** 1.22.x (CI witnesses ubuntu/macos/windows on 1.22.x, plus ubuntu “stable”).

## Book 2 suite map

This repo is designed to be used alongside the other Book 2 repos:

- **[finance-pipeline-gcp](https://github.com/nicholaskarlson/finance-pipeline-gcp)** — anchor drop-folder workflow (trigger → run → artifacts → markers)
- **[proof-first-event-contracts](https://github.com/nicholaskarlson/proof-first-event-contracts)** — event parsing contract + fixtures/goldens + expected-fail
- **[proof-first-deploy-gcp](https://github.com/nicholaskarlson/proof-first-deploy-gcp)** — deterministic deploy evidence (render + verify) + fixtures/goldens
- **[proof-first-casefiles](https://github.com/nicholaskarlson/proof-first-casefiles)** — engagement kits you can hand to a client (or use in teaching)

## Quickstart

Run the proof gate:

```bash
make verify
# (optional) Equivalent, if you want to run it directly:
# go test -count=1 ./...
```

Run the deterministic fixture demo (recomputes outputs and diffs against `fixtures/expected/**`):

```bash
go run ./cmd/eventcontracts demo --out ./out
```

## Fixture layout (proof gate)

Each fixture case is a folder name (the book references only these names).

Inputs live under:

- `fixtures/input/<case>/event.json` (CloudEvent / Eventarc envelope fixture)

Goldens live under:

- `fixtures/expected/<case>/decision.json` + `object_ref.json` **OR**
- `fixtures/expected/<case>/error.txt` (expected-fail cases)

## Output artifacts

On success, this repo emits:

- `decision.json` — stable “run vs ignore” decision (pretty JSON + trailing newline)
- `object_ref.json` — stable parsed object identity (pretty JSON + trailing newline)

For expected-fail (contract violations), the output is:

- `error.txt` only (must end with a newline)

## Rules enforced (book-friendly, audit-friendly)

- **Deterministic ordering** and stable JSON formatting.
- **No entropy** in artifacts: no timestamps, UUIDs, random IDs, or host-specific paths.
- **LF-only** output text artifacts (repo normalized via `.gitattributes`).
- **Expected-fail** cases produce **only** `error.txt`.

## Docs

- `docs/CONVENTIONS.md` — determinism + expected-fail rules (suite tone)
- `docs/CONTRACT.md` — what inputs are accepted and what artifacts are emitted
- `docs/HANDOFF.md` — how downstream services should use these artifacts

## License

MIT (see `LICENSE`).
