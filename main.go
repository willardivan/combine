package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// DirNode represents a node in the directory tree
type DirNode struct {
	Name     string
	IsDir    bool
	Children []*DirNode
}

// buildDirTree builds a directory tree structure
func buildDirTree(rootPath string, outFilePath string, isTextOnly bool, includeExts, excludeExts []string) (*DirNode, error) {
	outAbs, err := filepath.Abs(outFilePath)
	if err != nil {
		return nil, err
	}

	rootNode := &DirNode{
		Name:     filepath.Base(rootPath),
		IsDir:    true,
		Children: []*DirNode{},
	}
	
	nodesMap := make(map[string]*DirNode)
	nodesMap[rootPath] = rootNode

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the output file itself
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if pathAbs == outAbs {
			return nil
		}

		// Skip the root path
		if path == rootPath {
			return nil
		}
		
		// Skip files that don't match extension filters
		if !info.IsDir() {
			ext := strings.TrimPrefix(filepath.Ext(path), ".")
			if len(includeExts) > 0 && !containsExt(includeExts, ext) {
				return nil
			}
			if len(excludeExts) > 0 && containsExt(excludeExts, ext) {
				return nil
			}
		}
		
		// Skip non-text files if isTextOnly is true
		if !info.IsDir() && isTextOnly && !isTextFile(path) {
			return nil
		}

		parentPath := filepath.Dir(path)
		parentNode, exists := nodesMap[parentPath]
		if !exists {
			return fmt.Errorf("parent node not found for %s", path)
		}

		node := &DirNode{
			Name:     filepath.Base(path),
			IsDir:    info.IsDir(),
			Children: []*DirNode{},
		}
		
		parentNode.Children = append(parentNode.Children, node)
		if info.IsDir() {
			nodesMap[path] = node
		}
		
		return nil
	})

	// Sort children alphabetically with directories first
	for _, node := range nodesMap {
		sort.Slice(node.Children, func(i, j int) bool {
			if node.Children[i].IsDir != node.Children[j].IsDir {
				return node.Children[i].IsDir
			}
			return node.Children[i].Name < node.Children[j].Name
		})
	}

	return rootNode, err
}

// containsExt checks if an extension is in the given list
func containsExt(exts []string, ext string) bool {
	ext = strings.ToLower(ext)
	for _, e := range exts {
		if strings.ToLower(e) == ext {
			return true
		}
	}
	return false
}

// getFormatStats returns statistics about file formats in the directory
func getFormatStats(rootPath string) (map[string]int, error) {
	stats := make(map[string]int)
	
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Get file extension (without the dot)
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
		if ext == "" {
			ext = "no-extension"
		}
		
		stats[ext]++
		return nil
	})
	
	return stats, err
}

// printTreeToString renders the tree structure to a string
func printTreeToString(node *DirNode, prefix string, isLast bool, result *strings.Builder) {
	if node.Name == "." || node.Name == "" {
		result.WriteString("Directory Structure:\n")
	} else {
		// Print current node
		entry := prefix
		if isLast {
			entry += "└── "
			prefix += "    "
		} else {
			entry += "├── "
			prefix += "│   "
		}
		
		result.WriteString(entry + node.Name)
		if node.IsDir {
			result.WriteString("/")
		}
		result.WriteString("\n")
	}

	// Print children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		printTreeToString(child, prefix, isLastChild, result)
	}
}

func main() {
	// Customize usage information.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output_file] [-f extensions] [-fe excluded_extensions] [-checkformat] [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCombine all text files in a directory (recursively) into a single output file with headers.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  directory     The directory to scan. If omitted, the current directory is used after confirmation.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Define flags
	outName := flag.String("o", "combined_text.txt", "output file name (default \"combined_text.txt\")")
	includeFormats := flag.String("f", "", "only include files with these extensions (comma-separated, e.g. \"py,txt,json\")")
	excludeFormats := flag.String("fe", "", "exclude files with these extensions (comma-separated, e.g. \"exe,jpg,png\")")
	checkFormatFlag := flag.Bool("checkformat", false, "check and display statistics about file formats in the directory")
	
	// Parse flags
	flag.Parse()
	
	// Determine the target directory.
	directory := "."
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	}

	// Handle --checkformat flag first, before any other operations
	if *checkFormatFlag {
		stats, err := getFormatStats(directory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking formats: %v\n", err)
			os.Exit(1)
		}
		
		// Create a sorted list of extensions for output
		extensions := make([]string, 0, len(stats))
		for ext := range stats {
			extensions = append(extensions, ext)
		}
		sort.Slice(extensions, func(i, j int) bool {
			// Sort by count (descending)
			if stats[extensions[i]] != stats[extensions[j]] {
				return stats[extensions[i]] > stats[extensions[j]]
			}
			// If counts are equal, sort alphabetically
			return extensions[i] < extensions[j]
		})
		
		// Display the results
		fmt.Printf("File format statistics for %s:\n", directory)
		fmt.Println("------------------------------------")
		total := 0
		for _, ext := range extensions {
			count := stats[ext]
			total += count
			fmt.Printf("%5d %s files\n", count, ext)
		}
		fmt.Println("------------------------------------")
		fmt.Printf("Total: %d files\n", total)
		os.Exit(0)
	}

	// Prompt for confirmation if no directory argument is provided.
	if flag.NArg() == 0 {
		fmt.Print("Are you sure you want to combine the current folder? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled.")
			os.Exit(0)
		}
	}
	
	// Process include/exclude extensions
	var includeExts, excludeExts []string
	if *includeFormats != "" {
		includeExts = strings.Split(*includeFormats, ",")
	}
	if *excludeFormats != "" {
		excludeExts = strings.Split(*excludeFormats, ",")
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

	// Build and write the directory tree structure
	dirTree, err := buildDirTree(directory, *outName, true, includeExts, excludeExts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building directory tree: %v\n", err)
		os.Exit(1)
	}
	
	// Convert tree to string
	treeBuilder := &strings.Builder{}
	printTreeToString(dirTree, "", false, treeBuilder)
	treeStr := treeBuilder.String()
	
	// Write tree structure to output file
	if _, err := outFile.WriteString(treeStr); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing directory structure: %v\n", err)
		os.Exit(1)
	}
	
	// Add filter information if filters are applied
	if len(includeExts) > 0 || len(excludeExts) > 0 {
		outFile.WriteString("\nFilters applied:\n")
		if len(includeExts) > 0 {
			outFile.WriteString(fmt.Sprintf("- Including only: %s\n", strings.Join(includeExts, ", ")))
		}
		if len(excludeExts) > 0 {
			outFile.WriteString(fmt.Sprintf("- Excluding: %s\n", strings.Join(excludeExts, ", ")))
		}
	}
	
	// Add separator between structure and content
	if _, err := outFile.WriteString("\n" + strings.Repeat("-", 80) + "\n\n"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing separator: %v\n", err)
		os.Exit(1)
	}

	// Walk the directory recursively.
	fileCount := 0
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

		// Apply extension filters
		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		ext = strings.ToLower(ext)
		if len(includeExts) > 0 && !containsExt(includeExts, ext) {
			return nil
		}
		if len(excludeExts) > 0 && containsExt(excludeExts, ext) {
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
		fileCount++
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(*outName)
	fmt.Printf("\nMerging complete. Output file: %s (%d files processed)\n", *outName, fileCount)
}
