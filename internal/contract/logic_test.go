package contract

import "testing"

func TestParseAndDecide_FinalizedGood(t *testing.T) {
	raw := []byte(`{
  "specversion": "1.0",
  "type": "google.cloud.storage.object.v1.finalized",
  "source": "//storage.googleapis.com/projects/_/buckets/pf-drop-bucket",
  "id": "evt-1",
  "data": { "bucket": "pf-drop-bucket", "name": "drop/left.csv" }
}`)
	dec, ref, errText := ParseAndDecide(raw, "pf-drop-bucket")
	if errText != nil {
		t.Fatalf("unexpected error: %q", *errText)
	}
	if !dec.ShouldRun || dec.Decision != "run" {
		t.Fatalf("expected run, got: %+v", dec)
	}
	if ref.NameUnescaped != "drop/left.csv" {
		t.Fatalf("unexpected unescaped name: %q", ref.NameUnescaped)
	}
}

func TestParseAndDecide_URLUnescape(t *testing.T) {
	raw := []byte(`{
  "specversion": "1.0",
  "type": "google.cloud.storage.object.v1.finalized",
  "id": "evt-2",
  "data": { "bucket": "pf-drop-bucket", "name": "drop%2Fright.csv" }
}`)
	dec, ref, errText := ParseAndDecide(raw, "pf-drop-bucket")
	if errText != nil {
		t.Fatalf("unexpected error: %q", *errText)
	}
	if !dec.ShouldRun {
		t.Fatalf("expected should_run true, got: %+v", dec)
	}
	if ref.Name != "drop%2Fright.csv" {
		t.Fatalf("expected raw name preserved, got %q", ref.Name)
	}
	if ref.NameUnescaped != "drop/right.csv" {
		t.Fatalf("expected unescaped, got %q", ref.NameUnescaped)
	}
}

func TestParseAndDecide_IgnoreNonFinalize(t *testing.T) {
	raw := []byte(`{
  "specversion": "1.0",
  "type": "google.cloud.storage.object.v1.deleted",
  "id": "evt-3",
  "data": { "bucket": "pf-drop-bucket", "name": "drop/left.csv" }
}`)
	dec, _, errText := ParseAndDecide(raw, "pf-drop-bucket")
	if errText != nil {
		t.Fatalf("unexpected error: %q", *errText)
	}
	if dec.ShouldRun {
		t.Fatalf("expected ignore, got: %+v", dec)
	}
}

func TestParseAndDecide_MissingRequiredFields(t *testing.T) {
	raw := []byte(`{
  "specversion": "1.0",
  "type": "google.cloud.storage.object.v1.finalized",
  "id": "evt-4",
  "data": { "bucket": "pf-drop-bucket" }
}`)
	_, _, errText := ParseAndDecide(raw, "pf-drop-bucket")
	if errText == nil {
		t.Fatalf("expected error")
	}
}
