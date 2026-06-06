package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	var (
		input      = flag.String("in", "", "Path to a password-protected PDF or a directory of PDFs (required)")
		outputFile = flag.String("out", "", "Output PDF path — ignored when -in is a directory")
		password   = flag.String("password", "", "PDF password (required)")
	)
	flag.Parse()

	if *input == "" || *password == "" {
		fmt.Println("Usage: go run unlocker.go -in=file.pdf -password=secret [-out=output.pdf]")
		fmt.Println("       go run unlocker.go -in=directory/  -password=secret")
		flag.PrintDefaults()
		os.Exit(1)
	}

	info, err := os.Stat(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", *input, err)
		os.Exit(1)
	}

	if info.IsDir() {
		if err := unlockDir(*input, *password); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		out := *outputFile
		if out == "" {
			out = defaultOutputPath(*input)
		}
		if err := unlockFile(*input, out, *password); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func unlockDir(dir, password string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	var pdfs []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".pdf") {
			pdfs = append(pdfs, filepath.Join(dir, e.Name()))
		}
	}

	if len(pdfs) == 0 {
		fmt.Printf("No PDFs found in: %s\n", dir)
		return nil
	}

	fmt.Printf("Found %d PDF(s) in: %s\n\n", len(pdfs), dir)

	ok, failed := 0, 0
	for _, path := range pdfs {
		out := defaultOutputPath(path)
		if err := unlockFile(path, out, password); err != nil {
			fmt.Fprintf(os.Stderr, "  ERROR: %v\n", err)
			failed++
		} else {
			ok++
		}
	}

	fmt.Printf("\nDone: %d OK, %d failed\n", ok, failed)
	return nil
}

func unlockFile(input, output, password string) error {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password

	fmt.Printf("Unlocking: %s → %s\n", input, output)

	if err := api.DecryptFile(input, output, conf); err != nil {
		return fmt.Errorf("removing password from %s: %w", filepath.Base(input), err)
	}

	fmt.Printf("  OK\n")
	return nil
}

func defaultOutputPath(input string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	return filepath.Join(dir, name+"_unlocked"+ext)
}
