package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// isTextFile checks whether a file appears to be a text file.
// It reads the first 512 bytes and returns false if a null byte is found
// or if the data is not valid UTF-8.
func isTextFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}
	// Check for null bytes.
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return false
		}
	}
	// Check for valid UTF-8.
	if !utf8.Valid(buf[:n]) {
		return false
	}
	return true
}

func main() {
	// Customize usage information.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [directory] [-o output_file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCombine all text files in a directory (recursively) into a single output file with headers.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  directory     The directory to scan. If omitted, the current directory is used after confirmation.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Define the output file flag.
	outName := flag.String("o", "combined_text.txt", "output file name (default \"combined_text.txt\")")
	flag.Parse()

	// Determine the target directory.
	directory := "."
	if flag.NArg() == 0 {
		// Prompt for confirmation if no directory argument is provided.
		fmt.Print("Are you sure you want to combine the current folder? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}
	} else {
		directory = flag.Arg(0)
	}

	// Create (or truncate) the output file.
	outFile, err := os.Create(*outName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Get the absolute path of the output file to avoid processing it.
	outAbs, err := filepath.Abs(*outName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error obtaining absolute path of output file: %v\n", err)
		os.Exit(1)
	}

	// Walk the directory recursively.
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Skip the output file itself.
		currAbs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if currAbs == outAbs {
			return nil
		}

		// Only process text files.
		if !isTextFile(path) {
			return nil
		}

		// Get the relative path for the header.
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			relPath = path
		}

		// Write the header.
		header := fmt.Sprintf("== %s ==\n", relPath)
		if _, err := outFile.WriteString(header); err != nil {
			return err
		}

		// Read and write the file content.
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := outFile.Write(content); err != nil {
			return err
		}

		// Add spacing between file contents.
		if _, err := outFile.WriteString("\n\n"); err != nil {
			return err
		}

		// Output the processed file name to the console.
		fmt.Println(relPath)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(*outName)
	fmt.Printf("\nMerging complete. Output file: %s\n", *outName)
}
