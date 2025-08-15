# Contributing to qage

Thank you for your interest in contributing to qage! This document provides guidelines for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/zlobste/qage.git
   cd qage
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   go build -o bin/qage ./cmd/qage
   go build -o bin/age-plugin-qage ./cmd/age-plugin-qage
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## Code Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `golangci-lint` to check for issues
- Add comments for exported functions and types

### Package Structure

- **`pkg/`** - Public APIs that can be imported by other projects
- **`internal/`** - Private packages not meant for external use
- **`cmd/`** - Executable applications
- **`docs/`** - Generated documentation

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass: `go test ./...`
- Run built-in validation: `./bin/qage selftest`
- Test CLI commands manually

## Submitting Changes

### Pull Request Process

1. **Fork** the repository on GitHub
2. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following the code guidelines
4. **Add tests** for new functionality
5. **Run the test suite**:
   ```bash
   go test ./...
   go vet ./...
   golangci-lint run
   ```
6. **Commit your changes** with clear messages:
   ```bash
   git commit -m "Add feature: brief description"
   ```
7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
8. **Create a Pull Request** on GitHub

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add support for custom key comments

- Allow users to add comments when generating keys
- Update CLI help text and documentation
- Add tests for comment functionality

Fixes #123
```

## Development Workflow

### Before Submitting

1. **Run all checks**:
   ```bash
   # Format code
   gofmt -w .
   
   # Run tests
   go test ./...
   
   # Check for issues
   go vet ./...
   golangci-lint run
   
   # Test CLI
   ./bin/qage selftest
   ```

2. **Update documentation** if needed:
   ```bash
   # Regenerate CLI docs
   ./bin/qage docs
   
   # Check that docs/ directory is updated
   git status docs/
   ```

3. **Check that examples work**:
   ```bash
   go run examples/library_usage.go
   ```

### Continuous Integration

Our CI pipeline runs:
- Tests on Go 1.24+
- Linting with golangci-lint
- Security scanning
- Cross-platform builds

Make sure your changes pass all CI checks before submitting.

## Types of Contributions

### Bug Reports

- Use the GitHub issue tracker
- Include steps to reproduce
- Provide system information (Go version, OS)
- Include relevant error messages

### Feature Requests

- Check existing issues first
- Describe the use case
- Explain why this feature would be useful
- Consider backward compatibility

### Code Contributions

We welcome:
- Bug fixes
- Performance improvements
- New CLI commands
- Library enhancements
- Documentation improvements
- Test coverage improvements

## Security

For security-related issues:
- **Do not** open a public issue
- Email the maintainers directly
- Allow time for fixes before disclosure

## Questions?

- Open an issue for general questions
- Check existing documentation first
- Be respectful and constructive

Thank you for contributing to qage! 🚀
