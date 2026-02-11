package contract

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const (
	TypeFinalized       = "google.cloud.storage.object.v1.finalized"
	TypeDeleted         = "google.cloud.storage.object.v1.deleted"
	TypeMetadataUpdated = "google.cloud.storage.object.v1.metadataUpdated"
)

// ObjectRef is the parsed bucket + object name.
// NameUnescaped is url.PathUnescape(Name) when possible.
type ObjectRef struct {
	Bucket        string `json:"bucket"`
	Name          string `json:"name"`
	NameUnescaped string `json:"name_unescaped"`
}

// Decision is the stable, reviewable “run vs ignore” outcome.
type Decision struct {
	ShouldRun bool   `json:"should_run"`
	Decision  string `json:"decision"` // "run" or "ignore"
	Reason    string `json:"reason"`

	EventType string `json:"event_type"`
	EventID   string `json:"event_id,omitempty"`

	Bucket          string `json:"bucket,omitempty"`
	Object          string `json:"object,omitempty"`
	ObjectUnescaped string `json:"object_unescaped,omitempty"`
}

// ParseAndDecide parses a CloudEvent and returns deterministic artifacts.
// If the event cannot be parsed safely, it returns a deterministic error string.
func ParseAndDecide(raw []byte, expectedBucket string) (Decision, ObjectRef, *string) {
	// 1) Parse CloudEvent JSON (strict enough to catch malformed payloads).
	var ev CloudEvent
	if err := json.Unmarshal(raw, &ev); err != nil {
		msg := "error: invalid JSON (failed to parse CloudEvent)\n"
		return Decision{}, ObjectRef{}, &msg
	}
	if strings.TrimSpace(ev.SpecVersion) == "" || strings.TrimSpace(ev.Type) == "" {
		msg := "error: missing required CloudEvent fields (specversion/type)\n"
		return Decision{}, ObjectRef{}, &msg
	}

	// 2) Parse StorageData (required for any decision that references bucket/name).
	var data StorageData
	if len(ev.Data) > 0 {
		_ = json.Unmarshal(ev.Data, &data) // validate below
	}

	// Non-finalize events are ignored, but we still try to parse bucket/name if present.
	switch ev.Type {
	case TypeFinalized:
		// For finalized, bucket+name are required.
		if strings.TrimSpace(data.Bucket) == "" || strings.TrimSpace(data.Name) == "" {
			msg := "error: finalized event missing required fields (data.bucket/data.name)\n"
			return Decision{}, ObjectRef{}, &msg
		}
	default:
		// ignore (noise) — no additional required fields
	}

	bucket := clean(data.Bucket)
	name := clean(data.Name)

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

	// 3) Decision.
	dec := Decision{
		EventType: ev.Type,
		EventID:   clean(ev.ID),
	}

	if ev.Type != TypeFinalized {
		dec.ShouldRun = false
		dec.Decision = "ignore"
		dec.Reason = fmt.Sprintf("ignore: non-finalize event (%s)", ev.Type)
		dec.Bucket = ref.Bucket
		dec.Object = ref.Name
		dec.ObjectUnescaped = ref.NameUnescaped
		return dec, ref, nil
	}

	// Finalized.
	dec.Bucket = ref.Bucket
	dec.Object = ref.Name
	dec.ObjectUnescaped = ref.NameUnescaped

	if ref.Bucket != expectedBucket {
		dec.ShouldRun = false
		dec.Decision = "ignore"
		dec.Reason = fmt.Sprintf("ignore: bucket mismatch (got %q, want %q)", ref.Bucket, expectedBucket)
		return dec, ref, nil
	}

	dec.ShouldRun = true
	dec.Decision = "run"
	dec.Reason = "run: finalized event in expected bucket"
	return dec, ref, nil
}

func clean(s string) string {
	return strings.TrimSpace(s)
}
