# Contributing to Combine Tool

Thank you for your interest in contributing to Combine Tool! We welcome contributions from everyone. This document provides guidelines and instructions for contributing.

## Ways to Contribute

There are several ways you can contribute to this project:

- **Bug Reports**: Report bugs by opening an issue
- **Feature Requests**: Suggest new features or improvements
- **Documentation**: Improve or add documentation
- **Code Contributions**: Submit pull requests with bug fixes or new features

## Development Workflow

1. **Fork the repository**
2. **Clone your fork**
   ```
   git clone https://github.com/YOUR_USERNAME/combine.git
   cd combine
   ```
3. **Create a new branch**
   ```
   git checkout -b feature/your-feature-name
   ```
   or
   ```
   git checkout -b fix/your-bug-fix
   ```
4. **Make your changes**
5. **Test your changes**
   ```
   go build
   ./combine -v  # Verify version
   # Test with various options
   ```
6. **Commit your changes**
   ```
   git commit -m "Your descriptive commit message"
   ```
7. **Push to your fork**
   ```
   git push origin feature/your-feature-name
   ```
8. **Open a Pull Request**

## Code Style

- Follow standard Go code style and conventions
- Run `go fmt` before committing to ensure consistent formatting
- Write descriptive commit messages

## Reporting Bugs

When reporting bugs, please include:

- Detailed steps to reproduce the bug
- Expected behavior and actual behavior
- Environment details (OS, Go version, etc.)
- Any error messages or logs

## Feature Requests

When requesting new features, please:

- Clearly describe the feature and its benefits
- Explain how it would be used
- Consider potential implementation approaches if possible

## Pull Request Process

1. Update documentation if necessary
2. Make sure code compiles and runs without errors
3. Add or update tests if applicable
4. Your PR will be reviewed by maintainers, who may request changes
5. Once approved, your PR will be merged

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](LICENSE.md).

Thank you for your contributions! 