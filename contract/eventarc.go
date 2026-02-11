package contract

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ParseEventarcAndDecide parses an Eventarc-delivered event (Cloud Run) and returns deterministic artifacts.
//
// Inputs:
//   - ceType: Eventarc Ce-Type header value (event type)
//   - body: request body JSON (either direct {bucket,name} or envelope {data:{bucket,name}})
//   - expectedBucket: bucket name we will RUN for (finalized only)
//
// Outputs:
//   - Decision + ObjectRef (when parsed) OR stable error text (expected-fail contract).
func ParseEventarcAndDecide(ceType string, body []byte, expectedBucket string) (Decision, ObjectRef, *string) {
	typ := clean(ceType)
	if typ == "" {
		msg := "error: missing required Eventarc header (Ce-Type)\n"
		return Decision{}, ObjectRef{}, &msg
	}

	// Parse body JSON (supports both shapes).
	// Shape A: {"bucket":"...","name":"..."}
	// Shape B: {"data":{"bucket":"...","name":"..."}}
	var direct StorageData
	if err := json.Unmarshal(body, &direct); err != nil {
		msg := "error: invalid JSON (failed to parse Eventarc body)\n"
		return Decision{}, ObjectRef{}, &msg
	}

	bucket := clean(direct.Bucket)
	name := clean(direct.Name)

	if bucket == "" && name == "" {
		var env struct {
			Data StorageData `json:"data"`
		}
		// body already valid JSON; still safe to unmarshal into envelope.
		_ = json.Unmarshal(body, &env)
		bucket = clean(env.Data.Bucket)
		name = clean(env.Data.Name)
	}

	// For finalized, bucket+name are required.
	if typ == TypeFinalized {
		if bucket == "" || name == "" {
			msg := "error: finalized event missing required fields (bucket/name)\n"
			return Decision{}, ObjectRef{}, &msg
		}
	}

	nameUnescaped := ""
	if name != "" {
		u, err := url.PathUnescape(name)
		if err != nil {
			msg := "error: invalid URL-escaped object name\n"
			return Decision{}, ObjectRef{}, &msg
		}
		nameUnescaped = u
	}

	ref := ObjectRef{
		Bucket:        bucket,
		Name:          name,
		NameUnescaped: nameUnescaped,
	}

	dec := Decision{
		EventType:       typ,
		Bucket:          ref.Bucket,
		Object:          ref.Name,
		ObjectUnescaped: ref.NameUnescaped,
	}

	// Non-finalize events are deterministic ignores.
	if typ != TypeFinalized {
		dec.ShouldRun = false
		dec.Decision = "ignore"
		dec.Reason = fmt.Sprintf("ignore: non-finalize event (%s)", typ)
		return dec, ref, nil
	}

	// Finalized: enforce expected bucket.
	exp := clean(expectedBucket)
	if ref.Bucket != exp {
		dec.ShouldRun = false
		dec.Decision = "ignore"
		dec.Reason = fmt.Sprintf("ignore: bucket mismatch (got %q, want %q)", ref.Bucket, exp)
		return dec, ref, nil
	}

	dec.ShouldRun = true
	dec.Decision = "run"
	dec.Reason = "run: finalized event in expected bucket"
	return dec, ref, nil
}
