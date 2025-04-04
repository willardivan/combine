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

// Version information
const (
	Version     = "1.0.0"
	BuildDate   = "2025-04-05"
	Description = "A utility to recursively combine text files into a single output file"
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
func buildDirTree(rootPath string, outFilePath string, isTextOnly bool, includeExts, excludeExts, excludePaths []string, textPattern string) (*DirNode, error) {
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

	// The first pass is only to identify files matching the pattern, if a pattern is specified
	var patternMatchedFiles map[string]bool
	if textPattern != "" {
		patternMatchedFiles = make(map[string]bool)
		
		// First walk to identify files with the pattern
		err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories and output file
			if info.IsDir() {
				// Check for excluded directories
				relPath, err := filepath.Rel(rootPath, path)
				if err != nil {
					relPath = path
				}
				
				if isExcludedPath(relPath, excludePaths) {
					return filepath.SkipDir
				}
				
				return nil
			}
			
			// Skip the output file
			pathAbs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if pathAbs == outAbs {
				return nil
			}
			
			// Get relative path for exclusion check
			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				relPath = path
			}
			
			// Skip excluded files
			if isExcludedPath(relPath, excludePaths) {
				return nil
			}
			
			// Apply extension filters
			if !info.IsDir() {
				ext := strings.TrimPrefix(filepath.Ext(path), ".")
				if len(includeExts) > 0 && !containsExt(includeExts, ext) {
					return nil
				}
				if len(excludeExts) > 0 && containsExt(excludeExts, ext) {
					return nil
				}
			}
			
			// Only process text files if required
			if isTextOnly && !isTextFile(path) {
				return nil
			}
			
			// Check for pattern match
			if fileContainsPattern(path, textPattern) {
				patternMatchedFiles[path] = true
			}
			
			return nil
		})
		
		if err != nil {
			return nil, err
		}
	}

	// Second pass to build the tree, only with files that contain the pattern
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
		
		// Skip excluded paths
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			relPath = path
		}
		
		if isExcludedPath(relPath, excludePaths) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		// For files, check if they match the pattern (if pattern is specified)
		if !info.IsDir() && textPattern != "" && !patternMatchedFiles[path] {
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

		// Add directory nodes even if no files match, to maintain directory structure
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
func getFormatStats(rootPath string, includeExts, excludeExts, excludePaths []string, textPattern string) (map[string]int, error) {
	stats := make(map[string]int)
	
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories, but check if they should be excluded first
		if info.IsDir() {
			// Get relative path for exclusion check
			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				relPath = path
			}
			
			// Skip excluded directories
			if isExcludedPath(relPath, excludePaths) {
				return filepath.SkipDir
			}
			
			return nil
		}
		
		// Get relative path for exclusion check
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			relPath = path
		}
		
		// Skip excluded files
		if isExcludedPath(relPath, excludePaths) {
			return nil
		}
		
		// Get file extension (without the dot)
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
		
		// Apply extension filters
		if len(includeExts) > 0 && !containsExt(includeExts, ext) {
			return nil
		}
		if len(excludeExts) > 0 && containsExt(excludeExts, ext) {
			return nil
		}
		
		// Check if the file contains the pattern if specified
		if textPattern != "" && !fileContainsPattern(path, textPattern) {
			return nil
		}
		
		// Use "no-extension" for files without extension
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

// isExcludedPath checks if a path matches any of the excluded paths
func isExcludedPath(path string, excludedPaths []string) bool {
	if len(excludedPaths) == 0 {
		return false
	}

	// Normalize path separators for consistent matching
	normalizedPath := filepath.ToSlash(path)
	
	for _, excludedPath := range excludedPaths {
		// Normalize excluded path
		normalizedExcludedPath := filepath.ToSlash(excludedPath)
		
		// Check for exact match
		if normalizedPath == normalizedExcludedPath {
			return true
		}
		
		// Check if this is a directory prefix match
		// e.g. "node_modules" should match "node_modules/anything"
		if strings.HasPrefix(normalizedPath, normalizedExcludedPath+"/") {
			return true
		}
		
		// Check for path matching with glob patterns
		matched, err := filepath.Match(normalizedExcludedPath, normalizedPath)
		if err == nil && matched {
			return true
		}
	}
	
	return false
}

// fileContainsPattern checks if a file contains the specified text pattern
func fileContainsPattern(path, pattern string) bool {
	if pattern == "" {
		return true // Always match if no pattern is specified
	}

	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	// Convert to string and check if pattern exists
	contentStr := string(content)
	return strings.Contains(contentStr, pattern)
}

func main() {
	// Customize usage information.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Combine %s - %s\n\n", Version, Description)
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output_file] [-f extensions] [-fe excluded_extensions] [-e excluded_paths] [-p pattern] [-nocompact] [-checkformat] [directory]\n", os.Args[0])
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
	excludePaths := flag.String("e", ".git", "exclude specific files or directories (comma-separated paths, e.g. \"node_modules,dist,temp.txt\") (default \".git\")")
	pattern := flag.String("p", "", "only include files containing this text pattern")
	checkFormatFlag := flag.Bool("checkformat", false, "check and display statistics about file formats in the directory")
	noCompactFlag := flag.Bool("nocompact", false, "don't compress file content to single line (default is to compress)")
	versionFlag := flag.Bool("v", false, "display version information")
	
	// Parse flags
	flag.Parse()
	
	// Handle version flag
	if *versionFlag {
		fmt.Printf("Combine %s\n", Version)
		fmt.Printf("Build date: %s\n", BuildDate)
		fmt.Printf("%s\n", Description)
		os.Exit(0)
	}
	
	// Determine the target directory.
	directory := "."
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	}

	// Handle --checkformat flag first, before any other operations
	if *checkFormatFlag {
		// Process include/exclude extensions
		var includeExts, excludeExts []string
		if *includeFormats != "" {
			includeExts = strings.Split(*includeFormats, ",")
		}
		if *excludeFormats != "" {
			excludeExts = strings.Split(*excludeFormats, ",")
		}

		// Process exclude paths
		var excludedPaths []string
		if *excludePaths != "" {
			excludedPaths = strings.Split(*excludePaths, ",")
			// Trim spaces
			for i := range excludedPaths {
				excludedPaths[i] = strings.TrimSpace(excludedPaths[i])
			}
		}
		
		// Get stats with filters applied
		stats, err := getFormatStats(directory, includeExts, excludeExts, excludedPaths, *pattern)
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
		
		// Print filter information if any filters are applied
		if len(includeExts) > 0 || len(excludeExts) > 0 || len(excludedPaths) > 0 {
			fmt.Println("Filters applied:")
			if len(includeExts) > 0 {
				fmt.Printf("- Including only: %s\n", strings.Join(includeExts, ", "))
			}
			if len(excludeExts) > 0 {
				fmt.Printf("- Excluding extensions: %s\n", strings.Join(excludeExts, ", "))
			}
			if len(excludedPaths) > 0 {
				fmt.Printf("- Excluding paths: %s\n", strings.Join(excludedPaths, ", "))
			}
			fmt.Println()
		}
		
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

	// Process exclude paths
	var excludedPaths []string
	if *excludePaths != "" {
		excludedPaths = strings.Split(*excludePaths, ",")
		// Trim spaces
		for i := range excludedPaths {
			excludedPaths[i] = strings.TrimSpace(excludedPaths[i])
		}
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
	dirTree, err := buildDirTree(directory, *outName, true, includeExts, excludeExts, excludedPaths, *pattern)
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
	if len(includeExts) > 0 || len(excludeExts) > 0 || len(excludedPaths) > 0 || *pattern != "" {
		outFile.WriteString("\nFilters applied:\n")
		if len(includeExts) > 0 {
			outFile.WriteString(fmt.Sprintf("- Including only: %s\n", strings.Join(includeExts, ", ")))
		}
		if len(excludeExts) > 0 {
			outFile.WriteString(fmt.Sprintf("- Excluding extensions: %s\n", strings.Join(excludeExts, ", ")))
		}
		if len(excludedPaths) > 0 {
			outFile.WriteString(fmt.Sprintf("- Excluding paths: %s\n", strings.Join(excludedPaths, ", ")))
		}
		if *pattern != "" {
			outFile.WriteString(fmt.Sprintf("- Only files containing: \"%s\"\n", *pattern))
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
			// Check if directory is in excluded paths
			relPath, err := filepath.Rel(directory, path)
			if err != nil {
				relPath = path
			}
			
			if isExcludedPath(relPath, excludedPaths) {
				return filepath.SkipDir
			}
			
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

		// Get the relative path for checking exclusions
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			relPath = path
		}
		
		// Skip excluded files
		if isExcludedPath(relPath, excludedPaths) {
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
		
		// Check if file contains the specified pattern
		if *pattern != "" && !fileContainsPattern(path, *pattern) {
			return nil
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

		// If compact flag is set, compress content to a single line
		if !*noCompactFlag {
			// Replace newlines with a special delimiter that helps preserve code structure
			contentStr := string(content)
			
			// Normalize line endings
			contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
			contentStr = strings.ReplaceAll(contentStr, "\r", "\n")
			
			// Process each line and add an indicator of indentation level
			lines := strings.Split(contentStr, "\n")
			var compressed strings.Builder
			
			for _, line := range lines {
				// Count leading whitespace to preserve indentation info
				indent := 0
				for _, c := range line {
					if c == ' ' {
						indent++
					} else if c == '\t' {
						indent += 4 // Treat tab as 4 spaces
					} else {
						break
					}
				}
				
				// Trim the line
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == "" {
					continue // Skip empty lines
				}
				
				// Add a separator between lines, but not before the first line
				if compressed.Len() > 0 {
					compressed.WriteString(" ")
				}
				
				// Add indentation spaces for readability, without special symbols
				if indent > 0 {
					// Use a space followed by additional spaces for each level of indentation
					compressed.WriteString(strings.Repeat(" ", 1+(indent/4)))
				}
				
				// Add the line content
				compressed.WriteString(trimmedLine)
			}
			
			// Write the compressed content
			if _, err := outFile.WriteString(compressed.String()); err != nil {
				return err
			}
		} else {
			// Write the original content
			if _, err := outFile.Write(content); err != nil {
				return err
			}
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
