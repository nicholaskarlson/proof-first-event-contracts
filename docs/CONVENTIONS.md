# Conventions

These conventions exist so the repo can *prove* its outputs: deterministic artifacts + fixtures + goldens + a verification gate.

## Line endings
- **LF only** (enforced via `.gitattributes`).

## Determinism
Artifacts must be deterministic:
- stable ordering (sorted case lists; no filesystem walk order dependence)
- stable formatting (pretty JSON + trailing newline)
- no timestamps, UUIDs, random IDs, or other entropy in outputs

## Atomic writes
Artifacts are written via temp file → rename so partial files never appear.

## Expected-fail
Expected-fail cases emit **only**:
- `error.txt` (must end with a newline)

The proof gate (`make verify`) recomputes outputs and compares them byte-for-byte to `fixtures/expected/**`.
