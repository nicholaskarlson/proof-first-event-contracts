package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  eventcontracts demo --out <dir>")
	fmt.Fprintln(os.Stderr, "  eventcontracts parse --in <file> --out <dir> --bucket <expectedBucket> [--eventarc --ce-type <Ce-Type>]")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "demo":
		os.Exit(cmdDemo(os.Args[2:]))
	case "parse":
		os.Exit(cmdParse(os.Args[2:]))
	default:
		usage()
		os.Exit(2)
	}
}

func cmdParse(args []string) int {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	in := fs.String("in", "", "input JSON file")
	out := fs.String("out", "", "output directory")
	bucket := fs.String("bucket", "", "expected bucket (finalize run gate)")
	eventarc := fs.Bool("eventarc", false, "parse Eventarc-delivered body (Cloud Run) instead of CloudEvent")
	ceType := fs.String("ce-type", "", "Eventarc Ce-Type header value (used with --eventarc)")

	_ = fs.Parse(args)

	if *in == "" || *out == "" || *bucket == "" {
		usage()
		return 2
	}

	raw, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		return 1
	}

	if err := parseOne(raw, *out, *bucket, *eventarc, *ceType); err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		return 1
	}

	return 0
}
