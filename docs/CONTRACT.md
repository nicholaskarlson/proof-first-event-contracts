# Contract (Repo A)

This document defines the **stable, testable contract** for event parsing + filtering.

If you change any behavior here, treat it as a contract change:
update fixtures/goldens together and let CI be your witness.

---

## 1) Input

This repo accepts Cloud Storage / Eventarc-style CloudEvents as a fixture file:

- `event.json`

The fixture may be:
- a CloudEvent JSON representation, or
- an Eventarc envelope containing a CloudEvent payload

The contract is “best-effort parse with strict outputs”:
- if we can parse the event into a stable decision + object identity, we emit JSON artifacts
- if the fixture is malformed or missing required fields, we emit only `error.txt`

Additional strictness (book-friendly):
- unknown top-level fields in the CloudEvent envelope are rejected (expected-fail)

---

## 2) Output artifacts (success)

On success, a case produces:

- `decision.json`
- `object_ref.json`

Both are:
- UTF-8 JSON
- pretty-printed (stable indentation)
- terminated with a trailing newline

### decision.json

A stable decision of whether a downstream pipeline should run work.

Required fields:

- `should_run` (boolean)
- `decision` (`"run"` or `"ignore"`)
- `reason` (string; stable and reviewable)

`reason` must not include volatile content (timestamps, host paths, stack traces).

### object_ref.json

A stable object identity (when available):

- `bucket` (string)
- `name_raw` (string; as received)
- `name_unescaped` (string; URL-unescaped variant when applicable)

Downstream code should prefer `name_unescaped` for matching trigger rules.

---

## 3) Output artifacts (expected-fail)

On expected-fail, a case produces **only**:

- `error.txt`

Rules:

- one file only (no partial JSON artifacts)
- ends with exactly one newline
- content is stable and diff-friendly (no stack traces, no volatile paths)

---

## 4) Determinism requirements

- No timestamps, UUIDs, random IDs, or host-specific paths embedded into artifacts.
- Stable ordering for any lists that influence outputs.
- Stable JSON formatting with trailing newline.
