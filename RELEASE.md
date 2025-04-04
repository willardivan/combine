# Combine 1.0.0 Release

## Overview
Combine 1.0.0 is a powerful utility that recursively combines text files from a directory into a single output file with intelligent formatting. This initial release includes comprehensive filtering capabilities, smart formatting options, and file analysis features.

## Designed for AI/LLM Context
Combine was specifically designed to help developers prepare source code as input context for AI tools and Large Language Models. It addresses several key challenges when working with codebases and AI:

- **Context Window Optimization**: Compresses code to fit more content within limited AI context windows
- **Structure Preservation**: Maintains code organization with directory trees and indentation markers
- **Smart Filtering**: Targets only the most relevant files for your AI prompts
- **Balanced Token Usage**: Format preserves code readability while minimizing token consumption
- **Noise Reduction**: Automatically excludes binary files, irrelevant directories, and non-code elements

## Key Features

### File Operations
- **Directory Structure Display**: Visualizes the complete file structure at the top of output files
- **Recursive Scanning**: Automatically processes subdirectories to any depth
- **UTF-8 Detection**: Intelligently identifies and processes valid text files

### Advanced Filtering
- **Extension Filtering**: Include (`-f`) or exclude (`-fe`) files based on extensions
- **Path Exclusion**: Skip specific directories or files (`-e` flag, excludes `.git` by default)
- **Content Pattern Matching**: Find files containing specific text patterns with `-p` flag

### Output Formatting
- **Compact Mode (Default)**: Compresses files to single line with preserved indentation structure
- **Multiline Option**: Maintains original formatting when using `-nocompact` flag
- **File Statistics**: Generate reports on file types with the `-checkformat` flag

### Safety Features
- **Interactive Confirmation**: Prompts before processing current directory
- **Output File Protection**: Automatically excludes output file from processing
- **Relative Path Preservation**: Maintains original directory structure references

## Installation

```bash
# Using Go install
go install github.com/username/combine@latest

# From source
git clone https://github.com/username/combine.git
cd combine
go build
```

## Examples

Basic usage:
```bash
# Combine current directory files
combine -o combined.txt

# Specify a different directory
combine -o all_files.txt ./src
```

Content filtering:
```bash
# Only include Python files containing "def main"
combine -f py -p "def main" -o main_functions.txt ./src

# Exclude binary files and node_modules directory
combine -fe exe,dll,jpg,png -e node_modules -o clean.txt ./project
```

## Documentation
Full documentation is available in the [README.md](https://github.com/username/combine/blob/main/README.md) file.

## License
This release is available under the MIT License.

## Contributing
Contributions are welcome! Please see [CONTRIBUTING.md](https://github.com/username/combine/blob/main/CONTRIBUTING.md) for details. 