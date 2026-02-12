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


## Where this fits (Book 2)

This repo is designed to be used alongside the other Book 2 repos:

- Anchor: `finance-pipeline-gcp` — deployable drop-folder workflow (trigger → run → artifacts → markers)
- Repo A: `proof-first-event-contracts` — event parsing contract + fixtures/goldens + expected-fail
- Repo B: `proof-first-deploy-gcp` — deterministic deploy evidence (render + verify) + fixtures/goldens
- Repo C: `proof-first-casefiles` — engagement kits you can hand to a client (or use in teaching)

## Go baseline + CI witness

- Baseline: **Go 1.22.x**
- CI witness: **ubuntu/macos/windows** on Go 1.22.x, plus **stable** on ubuntu to catch drift.


## Layout (fixtures + goldens)

Fixtures (inputs):

Each case folder is one of two shapes:

1) CloudEvents JSON fixture:

- `fixtures/input/CASE/event.json`

2) Eventarc delivery fixture (Eventarc posts a body + the event type arrives in the `Ce-Type` header):

- `fixtures/input/CASE/body.json`
- `fixtures/input/CASE/ce_type.txt`

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

Eventarc-style fixture:

```bash
go run ./cmd/eventcontracts parse \
  --eventarc \
  --ce-type "$(cat fixtures/input/case16_eventarc_direct_finalized_run/ce_type.txt)" \
  --in fixtures/input/case16_eventarc_direct_finalized_run/body.json \
  --out ./tmp/case16 \
  --bucket pf-drop-bucket
cat ./tmp/case16/decision.json
cat ./tmp/case16/object_ref.json
```
