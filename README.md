COMBINE TOOL
============

A Go utility to recursively combine text files from a directory into a single output file.

FEATURES
--------
- Displays the directory structure as a tree at the beginning of the output file
- Scans directories recursively
- Detects text files (UTF-8 valid, no null bytes)
- Filter files by extension (include or exclude)
- Shows file format statistics for a directory
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
  combine [-o output.txt] [-f extensions] [-fe excluded_extensions] [-checkformat] [directory]

Examples:
  # Combine current directory (with confirmation)
  combine -o combined.txt

  # Specify directory and output
  combine -o all_text.md ./docs

  # Only include specific file types
  combine -f py,js,txt -o code_only.txt ./src

  # Exclude specific file types
  combine -fe exe,dll,jpg,png -o no_binaries.txt ./project

  # Check file formats in a directory
  combine -checkformat ./src

FLAGS
-----
  -o string           Output file name (default "combined_text.txt")
  -f string           Only include files with these extensions (comma-separated, e.g. "py,txt,json")
  -fe string          Exclude files with these extensions (comma-separated, e.g. "exe,jpg,png")
  -checkformat        Check and display statistics about file formats in the directory

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
- Excluding: exe, dll

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
