package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sters/yaml-diff/yamldiff"
)

func main() {
	file1 := flag.String("file1", "", "Target File 1 (diff says: -)")
	file2 := flag.String("file2", "", "Target File 2 (diff says: +)")
	flag.Parse()

	yamls1 := yamldiff.Load(load(*file1))
	yamls2 := yamldiff.Load(load(*file2))

	fmt.Printf("--- %s\n+++ %s\n\n", *file1, *file2)
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

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%+v, %s", err, f)
		return ""
	}

	return string(b)
}
