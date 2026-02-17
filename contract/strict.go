package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// decodeCloudEventStrict decodes a CloudEvent fixture with unknown-field rejection.
// This is a book-friendly contract: drift in the envelope shape is surfaced as an expected-fail.
func decodeCloudEventStrict(raw []byte) (CloudEvent, error) {
	var ev CloudEvent
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&ev); err != nil {
		return CloudEvent{}, err
	}
	// Reject trailing tokens (non-whitespace) after the first JSON value.
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return CloudEvent{}, errors.New("trailing data")
		}
		return CloudEvent{}, err
	}
	return ev, nil
}
