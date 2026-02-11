package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/nicholaskarlson/proof-first-event-contracts/contract"
	"github.com/nicholaskarlson/proof-first-event-contracts/internal/fsx"
)

const (
	fixturesIn  = "fixtures/input"
	fixturesExp = "fixtures/expected"
)

func demo(outDir string, expectedBucket string) error {
	// Clean outDir to avoid stale artifacts.
	if err := fsx.RemoveAll(outDir); err != nil {
		return fmt.Errorf("demo: remove out dir: %w", err)
	}
	if err := fsx.MustMkdirAll(outDir); err != nil {
		return fmt.Errorf("demo: mkdir out dir: %w", err)
	}

	cases, err := listCases(fixturesIn)
	if err != nil {
		return err
	}

	for _, c := range cases {
		inFile := filepath.Join(fixturesIn, c, "event.json")
		raw, err := os.ReadFile(inFile)
		if err != nil {
			return fmt.Errorf("demo: read fixture %s: %w", inFile, err)
		}

		dec, ref, errText := contract.ParseAndDecide(raw, expectedBucket)

		outCase := filepath.Join(outDir, c)
		if err := fsx.MustMkdirAll(outCase); err != nil {
			return fmt.Errorf("demo: mkdir case out dir: %w", err)
		}

		if errText != nil {
			if err := fsx.WriteText(filepath.Join(outCase, "error.txt"), *errText); err != nil {
				return fmt.Errorf("demo: write error.txt: %w", err)
			}
		} else {
			if err := fsx.WriteJSON(filepath.Join(outCase, "decision.json"), dec); err != nil {
				return fmt.Errorf("demo: write decision.json: %w", err)
			}
			if err := fsx.WriteJSON(filepath.Join(outCase, "object_ref.json"), ref); err != nil {
				return fmt.Errorf("demo: write object_ref.json: %w", err)
			}
		}

		expCase := filepath.Join(fixturesExp, c)
		if err := fsx.CompareDirs(expCase, outCase); err != nil {
			return fmt.Errorf("case %s: %w", c, err)
		}
	}

	return nil
}

func listCases(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("demo: read fixtures dir: %w", err)
	}
	var cases []string
	for _, e := range entries {
		if e.IsDir() {
			cases = append(cases, e.Name())
		}
	}
	sort.Strings(cases)
	return cases, nil
}
