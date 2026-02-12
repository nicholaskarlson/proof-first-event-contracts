package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func cmdDemo(args []string) int {
	outDir := "./out"
	if len(args) >= 2 && args[0] == "--out" {
		outDir = args[1]
	}

	_ = os.RemoveAll(outDir)
	_ = os.MkdirAll(outDir, 0o755)

	ents, err := os.ReadDir("fixtures/input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read fixtures: %v\n", err)
		return 1
	}

	var names []string
	for _, e := range ents {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	// Book fixtures are pinned to this bucket name.
	bucket := "pf-drop-bucket"

	for _, c := range names {
		caseDir := filepath.Join("fixtures/input", c)

		eventarc := false
		inPath := filepath.Join(caseDir, "event.json")
		if _, err := os.Stat(filepath.Join(caseDir, "body.json")); err == nil {
			eventarc = true
			inPath = filepath.Join(caseDir, "body.json")
		}

		raw, err := os.ReadFile(inPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", inPath, err)
			return 1
		}

		ceType := ""
		if eventarc {
			if b, err := os.ReadFile(filepath.Join(caseDir, "ce_type.txt")); err == nil {
				ceType = strings.TrimSpace(string(b))
			}
		}

		gotDir := filepath.Join(outDir, c)
		if err := parseOne(raw, gotDir, bucket, eventarc, ceType); err != nil {
			fmt.Fprintf(os.Stderr, "case %s: %v\n", c, err)
			return 1
		}

		expDir := filepath.Join("fixtures/expected", c)
		if same, why := dirsEqual(expDir, gotDir); !same {
			fmt.Fprintf(os.Stderr, "MISMATCH %s: %s\n", c, why)
			return 1
		}
	}

	fmt.Println("OK: demo outputs match fixtures (all case(s))")
	return 0
}

func dirsEqual(a, b string) (bool, string) {
	la, _ := os.ReadDir(a)
	lb, _ := os.ReadDir(b)

	var aa, bb []string
	for _, e := range la {
		if !e.IsDir() {
			aa = append(aa, e.Name())
		}
	}
	for _, e := range lb {
		if !e.IsDir() {
			bb = append(bb, e.Name())
		}
	}

	sort.Strings(aa)
	sort.Strings(bb)

	if strings.Join(aa, ",") != strings.Join(bb, ",") {
		return false, "file list differs"
	}

	for _, name := range aa {
		ba, _ := os.ReadFile(filepath.Join(a, name))
		bb, _ := os.ReadFile(filepath.Join(b, name))
		if !bytes.Equal(ba, bb) {
			return false, fmt.Sprintf("%s differs", name)
		}
	}

	return true, "ok"
}
