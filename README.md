# proof-first-event-contracts

Repo A for **Book 2** (*Proof-First Pipelines in the Cloud*): fixture-backed **event contracts** for Cloud Storage / Eventarc-style CloudEvents.

**Goal:** turn “incoming storage event” into **deterministic artifacts** you can review as diffs:
- `decision.json` — should we run or ignore?
- `object_ref.json` — parsed bucket + object name (with URL-unescaped variant)
- `error.txt` — expected-fail contract (stable error text, ends with a single newline)

This repo is designed to anchor Book 2 chapters while deployment repos evolve.

---

## Canonical commands (proof-first)

```bash
# Proof gate (one command)
make verify

# Portable (no Makefile)
go test -count=1 ./...
go run ./cmd/eventcontracts demo --out ./out
```

If this is green locally and in CI, the repo can:
1) recompute fixture outputs deterministically, and  
2) byte-verify them against goldens.

---

## Layout (fixtures + goldens)

Fixtures (inputs):

- `fixtures/input/CASE/event.json`

Goldens (expected outputs):

- `fixtures/expected/CASE/decision.json`
- `fixtures/expected/CASE/object_ref.json`
- `fixtures/expected/CASE/error.txt` (expected-fail only)

Demo outputs (recomputed):

- `out/CASE/...`

**Expected-fail rule:** if `fixtures/expected/CASE/error.txt` exists, the case must produce **only** `out/CASE/error.txt` and it must match byte-for-byte.

---

## What “decision” means

Current contract is intentionally small and matches the Book 2 workflow:

- **RUN** if the event is a *finalize* event **and** the bucket matches `--bucket` (default: `pf-drop-bucket`).
- **IGNORE** otherwise (non-finalize noise events, wrong bucket, etc.).
- **ERROR** (expected-fail) if the payload is invalid JSON or missing required fields for parsing.

Object names are URL-unescaped (e.g. `"drop%2Fleft.csv" → "drop/left.csv"`) before downstream use.

---

## Run one case manually (optional)

```bash
go run ./cmd/eventcontracts parse --in fixtures/input/case01_finalized_good/event.json --out ./tmp/case01
cat ./tmp/case01/decision.json
cat ./tmp/case01/object_ref.json
```

