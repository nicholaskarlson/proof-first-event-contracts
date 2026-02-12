package main

import (
	"os"
	"path/filepath"

	contract "github.com/nicholaskarlson/proof-first-event-contracts/contract"
)

func parseOne(raw []byte, outDir string, expectedBucket string, eventarc bool, ceType string) error {
	_ = os.MkdirAll(outDir, 0o755)

	var (
		dec     contract.Decision
		ref     contract.ObjectRef
		errText *string
	)

	if eventarc {
		dec, ref, errText = contract.ParseEventarcAndDecide(ceType, raw, expectedBucket)
	} else {
		dec, ref, errText = contract.ParseAndDecide(raw, expectedBucket)
	}

	if errText != nil {
		return os.WriteFile(filepath.Join(outDir, "error.txt"), []byte(*errText), 0o644)
	}

	if err := writeJSON(filepath.Join(outDir, "decision.json"), dec); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(outDir, "object_ref.json"), ref); err != nil {
		return err
	}
	return nil
}
