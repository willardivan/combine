COMBINE TOOL
============

A Go utility to recursively combine text files from a directory into a single output file.

FEATURES
--------
- Scans directories recursively
- Detects text files (UTF-8 valid, no null bytes)
- Excludes output file from processing
- Interactive confirmation for current directory
- Preserves relative paths in headers
- Shows real-time processing progress

INSTALLATION
------------
1. Install Go: https://golang.org/doc/install
2. Run:
   go install github.com/[YOUR_USERNAME]/combine@latest

USAGE
-----
Basic command:
  combine [directory] -o output.txt

Examples:
  # Combine current directory (with confirmation)
  combine -o combined.txt
  
  # Specify directory and output
  combine ./docs -o all_text.md
  
  # Use default output name
  combine ./src

FLAGS
-----
  -o string  Output file name (default "combined_text.txt")

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
== relative/path/file1.txt ==
[file contents]

== another/file2.txt ==
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