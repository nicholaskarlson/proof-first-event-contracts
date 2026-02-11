package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholaskarlson/proof-first-event-contracts/contract"
	"github.com/nicholaskarlson/proof-first-event-contracts/internal/fsx"
)

func parseOne(inFile, outDir, expectedBucket string) error {
	raw, err := os.ReadFile(inFile)
	if err != nil {
		return fmt.Errorf("parse: read input: %w", err)
	}

	dec, ref, errText := contract.ParseAndDecide(raw, expectedBucket)

	if err := fsx.MustMkdirAll(outDir); err != nil {
		return fmt.Errorf("parse: mkdir out dir: %w", err)
	}

	if errText != nil {
		return fsx.WriteText(filepath.Join(outDir, "error.txt"), *errText)
	}
	if err := fsx.WriteJSON(filepath.Join(outDir, "decision.json"), dec); err != nil {
		return err
	}
	if err := fsx.WriteJSON(filepath.Join(outDir, "object_ref.json"), ref); err != nil {
		return err
	}
	return nil
}
