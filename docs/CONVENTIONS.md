# Conventions (suite-wide)

This repo follows the same “proof-first” conventions as the rest of the toolkit:

- **Deterministic artifacts**: no timestamps, random IDs, or environment noise in verified outputs.
- **Stable ordering**: fixture cases are processed in sorted order.
- **LF line endings**: fixtures and goldens are committed with LF.
- **Atomic writes**: outputs are written via temp file → rename.
- **Expected-fail contract**: if `fixtures/expected/CASE/error.txt` exists, the case must produce only `out/CASE/error.txt` and it must match byte-for-byte (ending with a single newline).
