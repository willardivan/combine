COMBINE TOOL
============

A Go utility to recursively combine text files from a directory into a single output file.

FEATURES
--------
- Displays the directory structure as a tree at the beginning of the output file
- Scans directories recursively
- Detects text files (UTF-8 valid, no null bytes)
- Filter files by extension (include or exclude)
- Exclude specific files or directories from processing (excludes .git by default)
- Filter files by content with pattern matching
- Shows file format statistics for a directory (can be combined with filters)
- Compact mode by default to compress file content into a single line while preserving code structure
- Excludes output file from processing
- Interactive confirmation for current directory
- Preserves relative paths in headers
- Shows real-time processing progress

INSTALLATION
------------
1. Install Go: https://golang.org/doc/install
2. Run:
   go install github.com/willardivan/combine@latest

USAGE
-----
Basic command:
  combine [-o output.txt] [-f extensions] [-fe excluded_extensions] [-e excluded_paths] [-p pattern] [-nocompact] [-checkformat] [directory]

Examples:
  # Combine current directory (with confirmation)
  combine -o combined.txt

  # Specify directory and output
  combine -o all_text.md ./docs

  # Only include specific file types
  combine -f py,js,txt -o code_only.txt ./src

  # Exclude specific file types
  combine -fe exe,dll,jpg,png -o no_binaries.txt ./project

  # Exclude specific directories and files (adds to default .git)
  combine -e "node_modules,dist,.git,temp.txt" -o clean.txt ./project
  
  # Only include files containing a specific pattern
  combine -p "API_KEY" -o api_files.txt ./src
  
  # Combine all Python files containing "def main"
  combine -f py -p "def main" -o main_functions.txt ./src
  
  # Output in multi-line format (not compact)
  combine -nocompact -o readable.txt ./src

  # Check file formats in a directory
  combine -checkformat ./src
  
  # Check file formats with filters
  combine -checkformat -f py,js ./src
  
  # Check file formats excluding certain directories
  combine -checkformat -e "node_modules,.git" ./project
  
  # Check files containing a specific pattern
  combine -checkformat -p "import" ./src

FLAGS
-----
  -o string           Output file name (default "combined_text.txt")
  -f string           Only include files with these extensions (comma-separated, e.g. "py,txt,json")
  -fe string          Exclude files with these extensions (comma-separated, e.g. "exe,jpg,png")
  -e string           Exclude specific files or directories (comma-separated paths, e.g. "node_modules,dist,temp.txt") (default ".git")
  -p string           Only include files containing this text pattern
  -checkformat        Check and display statistics about file formats in the directory
  -nocompact          Don't compress file content to single line (default is to compress)

COMPACT MODE
------------
By default, the tool compresses each file to a single line:
- Preserves code structure using simple indentation with spaces
- Uses extra spaces to indicate indentation level (more spaces = deeper indentation)
- Skips empty lines for cleaner output
- Makes it easier for AI tools to process while keeping code structure visible

To disable compact mode and output the original multiline format, use the -nocompact flag.

Example of compact output:
```
def example(): x = 5  if x > 3:   print("Greater than 3")  else:   print("Not greater")
```

FILE CRITERIA
-------------
Valid text files must:
- Be valid UTF-8 encoded
- Contain no null bytes in first 512 bytes
Files skipped:
- Binary files
- Directories
- Output file itself
- Invalid UTF-8 files

OUTPUT FORMAT
-------------
# Directory structure displayed at the top
├── example/
│   ├── subfolder/
│   │   └── file2.txt
│   └── file1.txt

# Filter information (if applied)
Filters applied:
- Including only: py, js
- Excluding extensions: exe, dll
- Excluding paths: node_modules, .git, temp.txt
- Only files containing: "API_KEY"

--------------------------------------------------------------------------------

== example/file1.txt ==
[file contents]

== example/subfolder/file2.txt ==
[file contents]

POSSIBLE IMPROVEMENTS
---------------------
- Add file pattern matching
- Support custom headers
- Add encoding detection
- Implement dry-run mode
- Add file size limits

LICENSE
-------
Public Domain (Unlicense)
