# Combine

<p align="center">
  <img src="https://img.shields.io/badge/version-1.0.0-blue" alt="version"/>
  <img src="https://img.shields.io/badge/license-MIT-green" alt="license"/>
  <img src="https://img.shields.io/badge/go-%3E%3D%201.16-blue" alt="go version"/>
</p>

<p align="center">
  A powerful Go utility for recursively combining text files from a directory into a single output file with intelligent formatting.
  <br>
  <b>Perfect for creating AI/LLM context from source code repositories.</b>
</p>

---

## 🤖 Designed for AI/LLM Context

Combine was specifically designed to help developers prepare source code as input context for AI tools and Large Language Models. It addresses several key challenges when working with codebases and AI:

- **Context Window Optimization**: Compresses code to fit more content within limited AI context windows
- **Structure Preservation**: Maintains code organization with directory trees and indentation markers
- **Smart Filtering**: Targets only the most relevant files for your AI prompts
- **Balanced Token Usage**: Format preserves code readability while minimizing token consumption
- **Noise Reduction**: Automatically excludes binary files, irrelevant directories, and non-code elements

Perfect for:
- Asking LLMs to analyze specific parts of your codebase
- Creating comprehensive documentation from code comments
- Optimizing code review workflows with AI assistants
- Building custom AI training datasets from source repositories

## ✨ Features

- 📂 **Directory Structure Display** - Visualize the file structure at the top of the output
- 🔍 **Smart Filtering** - Filter by extension, content patterns, or specific paths
- 🚫 **Automatic Exclusions** - Skip binary files, .git directories (by default)
- 📊 **Format Statistics** - Analyze file types in your codebase with `-checkformat`
- 💼 **Content Pattern Matching** - Find files containing specific text patterns
- 📝 **Compact Mode** - Compress files to single line for easy AI processing (default)
- 🔄 **Multiline Option** - Keep original formatting with `-nocompact` when needed
- 🛡️ **Safety Measures** - Interactive confirmation, output file exclusion

## 📋 Installation

### Using Go Install

```bash
go install github.com/username/combine@latest
```

### From Source

```bash
git clone https://github.com/username/combine.git
cd combine
go build
```

## 🚀 Quick Start

### Basic Usage

```bash
# Combine current directory files (with confirmation)
combine -o combined.txt

# Specify a different directory
combine -o all_files.txt ./src
```

### Filtering Examples

```bash
# Only include Python, JavaScript, and text files
combine -f py,js,txt -o code_only.txt ./src

# Exclude binary and image files
combine -fe exe,dll,jpg,png -o no_binaries.txt ./project

# Exclude specific directories (adds to default .git)
combine -e "node_modules,dist,temp.txt" -o clean.txt ./project
```

### Content Pattern Matching

```bash
# Only include files containing API keys
combine -p "API_KEY" -o api_files.txt ./src

# Combine Python files containing "def main"
combine -f py -p "def main" -o main_functions.txt ./src
```

### Format Options

```bash
# Output in multi-line format (not compact)
combine -nocompact -o readable.txt ./src
```

### Analysis

```bash
# Check file formats in a directory
combine -checkformat ./src

# Check formats with filters
combine -checkformat -f py,js -e node_modules ./src
```

## 📖 Command-Line Options

| Flag | Description |
|------|-------------|
| `-o` | Output file name (default: "combined_text.txt") |
| `-f` | Include only these extensions (comma-separated, e.g., "py,txt,json") |
| `-fe` | Exclude these extensions (comma-separated, e.g., "exe,jpg,png") |
| `-e` | Exclude paths (comma-separated, e.g., "node_modules,dist") (default: ".git") |
| `-p` | Only include files containing this text pattern |
| `-checkformat` | Display statistics about file formats in the directory |
| `-nocompact` | Don't compress content to single line (default is to compress) |
| `-v` | Display version information |

## 💡 Compact Mode

By default, the tool compresses each file to a single line while preserving code structure:

```
def example(): x = 5  if x > 3:   print("Greater than 3")  else:   print("Not greater")
```

- Uses indentation with spaces to indicate nesting level
- Skips empty lines for cleaner output
- Perfect for AI tools and code analysis

Use `-nocompact` to disable this feature and output the original multiline format.

## 🧩 Output Format

```
├── example/
│   ├── subfolder/
│   │   └── file2.txt
│   └── file1.txt

Filters applied:
- Including only: py, js
- Excluding paths: node_modules, .git
- Only files containing: "API_KEY"

--------------------------------------------------------------------------------

== example/file1.txt ==
[file contents]

== example/subfolder/file2.txt ==
[file contents]
```

## 🔍 Advanced Use Cases

### AI Code Analysis Workflow

1. **Filter relevant code**:
   ```bash
   combine -f py,js -e "node_modules,tests" -o codebase.txt ./src
   ```

2. **Use in your AI prompt**:
   ```
   Analyze this codebase for security vulnerabilities:
   
   [Paste content of codebase.txt here]
   ```

### Finding Implementation Patterns

1. **Collect specific patterns**:
   ```bash
   combine -p "API.request" -o api_calls.txt ./src
   ```

2. **Ask AI to review**:
   ```
   Review all these API calls for error handling issues:
   
   [Paste content of api_calls.txt here]
   ```

### Repository Documentation

1. **Extract code with documentation comments**:
   ```bash
   combine -p "/**" -o docs.txt ./src
   ```

2. **Generate markdown documentation**:
   ```
   Convert these code comments to markdown documentation:
   
   [Paste content of docs.txt here]
   ```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md) for details on version updates and changes.

## 🔮 Roadmap

- [ ] Regular expression pattern matching
- [ ] Customizable section headers
- [ ] Better encoding detection for non-UTF8 files
- [ ] Dry-run mode
- [ ] File size limits for larger codebases
- [ ] Improved error handling and recovery
- [ ] Configuration file support
