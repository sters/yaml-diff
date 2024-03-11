package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AsakoKabe/yaml-diff/yamldiff"
)

func main() {
	ignoreEmptyFields := flag.Bool("ignore-empty-fields", false, "Ignore empty field")
	ignoreZeroFields := flag.Bool("ignore-zero-fields", false, "Ignore zero field")
	quiet := flag.Bool("quiet", false, "Print if diff exist")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: yaml-diff file1 file2")
		os.Exit(1)
	}
	file1 := args[0]
	file2 := args[1]

	yamls1, err := yamldiff.Load(load(file1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	yamls2, err := yamldiff.Load(load(file2))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
	}

	opts := []yamldiff.DoOptionFunc{}
	if *ignoreEmptyFields {
		opts = append(opts, yamldiff.EmptyAsNull())
	}
	if *ignoreZeroFields {
		opts = append(opts, yamldiff.ZeroAsNull())
	}

	if !*quiet {
		fmt.Printf("--- %s\n+++ %s\n\n", file1, file2)
	}
	for _, diff := range yamldiff.Do(yamls1, yamls2, opts...) {
		if *quiet && diff.Status() == yamldiff.DiffStatusSame {
			continue
		}
		fmt.Println(diff.Dump())
	}

	if !*quiet {
		fmt.Print()
	}
}

func load(f string) string {
	file, err := os.Open(f)
	defer func() { _ = file.Close() }()
	if err != nil {
		log.Printf("%+v, %s", err, f)

		return ""
	}

	b, err := io.ReadAll(file)
	if err != nil {
		log.Printf("%+v, %s", err, f)

		return ""
	}

	return string(b)
}
