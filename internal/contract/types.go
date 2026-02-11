package contract

import "encoding/json"

// CloudEvent is a minimal CloudEvents envelope for fixture contracts.
// We intentionally keep this struct small and strict.
type CloudEvent struct {
	SpecVersion string          `json:"specversion"`
	Type        string          `json:"type"`
	Source      string          `json:"source,omitempty"`
	ID          string          `json:"id,omitempty"`
	Subject     string          `json:"subject,omitempty"`
	Time        string          `json:"time,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// StorageData models the fields we care about for GCS object events.
type StorageData struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}
