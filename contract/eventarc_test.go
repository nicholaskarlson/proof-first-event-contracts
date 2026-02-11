package contract

import "testing"

func TestParseEventarcAndDecide_DirectShape_FinalizedRun(t *testing.T) {
	body := []byte(`{"bucket":"pf-drop-bucket","name":"drop/left.csv"}`)
	dec, ref, errText := ParseEventarcAndDecide(TypeFinalized, body, "pf-drop-bucket")
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

func TestParseEventarcAndDecide_EnvelopeShape_FinalizedRun(t *testing.T) {
	body := []byte(`{"data":{"bucket":"pf-drop-bucket","name":"drop%2Fright.csv"}}`)
	dec, ref, errText := ParseEventarcAndDecide(TypeFinalized, body, "pf-drop-bucket")
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

func TestParseEventarcAndDecide_MissingCeType_ExpectedFail(t *testing.T) {
	body := []byte(`{"bucket":"pf-drop-bucket","name":"drop/left.csv"}`)
	_, _, errText := ParseEventarcAndDecide("", body, "pf-drop-bucket")
	if errText == nil {
		t.Fatalf("expected error")
	}
}

func TestParseEventarcAndDecide_InvalidJSON_ExpectedFail(t *testing.T) {
	body := []byte(`{"bucket":`)
	_, _, errText := ParseEventarcAndDecide(TypeFinalized, body, "pf-drop-bucket")
	if errText == nil {
		t.Fatalf("expected error")
	}
}

func TestParseEventarcAndDecide_FinalizedMissingFields_ExpectedFail(t *testing.T) {
	body := []byte(`{"data":{"bucket":"pf-drop-bucket"}}`)
	_, _, errText := ParseEventarcAndDecide(TypeFinalized, body, "pf-drop-bucket")
	if errText == nil {
		t.Fatalf("expected error")
	}
}

func TestParseEventarcAndDecide_NonFinalizeIgnore(t *testing.T) {
	body := []byte(`{"bucket":"pf-drop-bucket","name":"drop/left.csv"}`)
	dec, _, errText := ParseEventarcAndDecide(TypeDeleted, body, "pf-drop-bucket")
	if errText != nil {
		t.Fatalf("unexpected error: %q", *errText)
	}
	if dec.ShouldRun {
		t.Fatalf("expected ignore, got: %+v", dec)
	}
}
