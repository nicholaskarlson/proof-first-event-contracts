package fsx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func MustMkdirAll(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func RemoveAll(dir string) error {
	return os.RemoveAll(dir)
}

// AtomicWriteFile writes content to path using a temp file in the same directory,
// then renames into place. It always writes with permissions 0644.
func AtomicWriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	// Best-effort cleanup.
	defer func() { _ = os.Remove(tmpName) }()

	if _, err := tmp.Write(content); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// Windows can't rename over an existing file.
	_ = os.Remove(path)
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	return nil
}

func WriteJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return AtomicWriteFile(path, b)
}

func WriteText(path, s string) error {
	if len(s) == 0 || s[len(s)-1] != '\n' {
		s += "\n"
	}
	return AtomicWriteFile(path, []byte(s))
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ListFiles(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)
	return files, nil
}

// CompareDirs compares the set of files in expectedDir to outDir and byte-compares each file.
// It rejects extra files in outDir (strict contract).
func CompareDirs(expectedDir, outDir string) error {
	expFiles, err := ListFiles(expectedDir)
	if err != nil {
		return fmt.Errorf("verify: list expected files: %w", err)
	}
	outFiles, err := ListFiles(outDir)
	if err != nil {
		return fmt.Errorf("verify: list output files: %w", err)
	}

	if !equalStringSlices(expFiles, outFiles) {
		return fmt.Errorf("verify: file set mismatch\nexpected: %v\n     got: %v", expFiles, outFiles)
	}

	for _, name := range expFiles {
		expPath := filepath.Join(expectedDir, name)
		outPath := filepath.Join(outDir, name)

		exp, err := os.ReadFile(expPath)
		if err != nil {
			return fmt.Errorf("verify: read expected %s: %w", expPath, err)
		}
		got, err := os.ReadFile(outPath)
		if err != nil {
			return fmt.Errorf("verify: read output %s: %w", outPath, err)
		}
		if !bytes.Equal(exp, got) {
			return fmt.Errorf("verify: golden mismatch: %s\nhint: diff -u %s %s", name, expPath, outPath)
		}
	}

	return nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// CopyReaderToFile is used rarely, but handy for deterministic test helpers.
func CopyReaderToFile(dst string, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return AtomicWriteFile(dst, b)
}
