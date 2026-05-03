// Command builddict parses the OpenCorpora XML dictionary and saves a
// compact binary dictionary in the .bin format for later loading by the
// gomorphy library.
//
// Usage:
//
//	builddict -in path/to/dict.opcorpora.xml -out dict.bin
//
// Filter: lemmas tagged NOUN (including the subtypes Name/Surn/Patr) and
// ADJF are kept. A "form → lemmas" reverse index is built in parallel for
// morphological analysis.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/therox/gomorphy/internal/dict"
)

func main() {
	var inPath, outPath string
	flag.StringVar(&inPath, "in", "", "path to the OpenCorpora XML dictionary")
	flag.StringVar(&outPath, "out", "", "path to the output .bin")
	flag.Parse()

	if inPath == "" || outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: builddict -in path.xml -out dict.bin")
		os.Exit(2)
	}

	if err := run(inPath, outPath); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(inPath, outPath string) error {
	in, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("opening XML: %w", err)
	}
	defer in.Close()

	d, stats, err := dict.BuildFromXML(in)
	if err != nil {
		return err
	}

	if err := dict.SaveFile(d, outPath); err != nil {
		return fmt.Errorf("saving dictionary: %w", err)
	}

	fi, err := os.Stat(outPath)
	if err != nil {
		return fmt.Errorf("stat of the output: %w", err)
	}

	fmt.Printf("lemmas processed: %d, kept: %d, forms: %d, tags: %d\n",
		stats.Seen, stats.Kept, stats.Forms, len(d.Tags))
	fmt.Printf("written to %s, size: %d bytes\n", outPath, fi.Size())
	return nil
}
