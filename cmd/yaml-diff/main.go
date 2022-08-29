package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sters/yaml-diff/yamldiff"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: yaml-diff file1 file2")
		os.Exit(1)
	}

	file1 := os.Args[1]
	file2 := os.Args[2]

	yamls1 := yamldiff.Load(load(file1))
	yamls2 := yamldiff.Load(load(file2))

	fmt.Printf("--- %s\n+++ %s\n\n", file1, file2)
	for _, diff := range yamldiff.Do(yamls1, yamls2) {
		fmt.Println(diff.Diff)
	}

	fmt.Print()
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
