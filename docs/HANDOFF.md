# Handoff (Repo A)

This document describes how downstream systems should use the event-contract artifacts.

---

## What to hand off

For a given event case, the handoff artifacts are:

- `decision.json` and `object_ref.json` (success), or
- `error.txt` only (expected-fail)

These artifacts are intentionally small so they can be:
- reviewed as PR diffs
- used as fixtures in other repos
- attached to an engagement handoff bundle

---

## How a pipeline should use the outputs

### If output is `error.txt`

- Treat the event as malformed/unsupported.
- ACK the event and do **not** run work.
- Do not retry based on this artifact alone.

### If output is `decision.json` + `object_ref.json`

- If `decision.json` says **ignore**, ACK the event and do **not** run work.
- If `decision.json` says **run**, apply your pipeline’s trigger rule using:
  - the bucket/object identity from `object_ref.json` (prefer `name_unescaped`)

---

## How this fits the Book 2 suite

- Repo A (**proof-first-event-contracts**) defines what events count as “run” vs “ignore” vs “error”.
- Repo 1 (**finance-pipeline-gcp**) consumes these outputs to implement replay-safe runs and completion markers.
- Repo B (**proof-first-deploy-gcp**) proves the deploy configuration itself (render + snapshot + verify).
- Repo C (**proof-first-casefiles**) provides realistic “engagement kits” that demonstrate the handoff loop.
