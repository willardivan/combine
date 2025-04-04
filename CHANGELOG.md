# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-04-05

### Added
- Initial release
- Directory structure visualization at the beginning of the output
- Filter files by extension (include/exclude)
- Exclude specific files or directories (default: .git)
- Pattern matching to find files containing specific text
- File format statistics with the -checkformat flag
- Compact mode by default (compresses to single line with indentation)
- Support for -nocompact flag to preserve original formatting
- Version information flag (-v)

### Changed
- Renamed output flag to -o

### Removed
- None 