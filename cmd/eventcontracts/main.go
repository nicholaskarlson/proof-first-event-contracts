package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	switch args[0] {
	case "demo":
		return cmdDemo(args[1:])
	case "parse":
		return cmdParse(args[1:])
	case "-h", "--help", "help":
		usage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		usage()
		return 2
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `proof-first-event-contracts

Commands:
  demo   Recompute fixtures into --out and byte-verify against goldens.
  parse  Parse a single event.json into --out (no golden verify).

Examples:
  make verify
  go run ./cmd/eventcontracts demo --out ./out
  go run ./cmd/eventcontracts parse --in fixtures/input/case01_finalized_good/event.json --out ./tmp/case01

`)
}

func cmdDemo(args []string) int {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	outDir := fs.String("out", "", "output directory (required)")
	bucket := fs.String("bucket", "pf-drop-bucket", "expected bucket name for RUN decisions")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *outDir == "" {
		fmt.Fprintln(os.Stderr, "--out is required")
		return 2
	}

	if err := demo(*outDir, *bucket); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	fmt.Println("OK: demo outputs match fixtures (all case(s))")
	return 0
}

func cmdParse(args []string) int {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	inFile := fs.String("in", "", "path to event JSON file (required)")
	outDir := fs.String("out", "", "output directory (required)")
	bucket := fs.String("bucket", "pf-drop-bucket", "expected bucket name for RUN decisions")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *inFile == "" || *outDir == "" {
		fmt.Fprintln(os.Stderr, "--in and --out are required")
		return 2
	}

	if err := parseOne(*inFile, *outDir, *bucket); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return 0
}
