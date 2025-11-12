# Contributing to OSS Insight Trending Notifier

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

1. A clear, descriptive title
2. Steps to reproduce the issue
3. Expected behavior
4. Actual behavior
5. Your environment (OS, Go version, etc.)
6. Any relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are welcome! Please create an issue with:

1. A clear description of the enhancement
2. The motivation/use case for this feature
3. Examples of how it would work
4. Any potential drawbacks or alternatives

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, concise commits
3. **Add tests** if applicable
4. **Update documentation** if needed
5. **Ensure tests pass**: `make test`
6. **Format your code**: `make fmt`
7. **Run linter**: `make vet`
8. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)
- Git

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/github-insight-analyze.git
cd github-insight-analyze
```

2. Install dependencies:
```bash
make deps
```

3. Set up configuration:
```bash
make setup-env
make setup-config
```

4. Build the project:
```bash
make build
```

5. Run tests:
```bash
make test
```

## Project Structure

```
.
├── cmd/notifier/          # Main application
├── pkg/                   # Public packages
│   ├── api/              # API client
│   ├── email/            # Email functionality
│   └── formatter/        # Data formatting
├── internal/config/       # Configuration management
├── configs/               # Configuration files
└── .github/workflows/     # CI/CD workflows
```

## Coding Guidelines

### Go Style

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Write comments for exported functions and types
- Keep functions focused and concise

### Commit Messages

Follow the conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(api): add support for custom API endpoints

fix(email): resolve SMTP authentication issue

docs(readme): update installation instructions
```

### Testing

- Write unit tests for new functionality
- Aim for good test coverage
- Use table-driven tests when appropriate
- Mock external dependencies

Example test structure:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    Type
        expected Type
        wantErr  bool
    }{
        {
            name:     "test case 1",
            input:    value1,
            expected: expected1,
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Error Handling

- Always handle errors explicitly
- Wrap errors with context using `fmt.Errorf` with `%w`
- Log errors appropriately
- Return meaningful error messages

Example:
```go
if err := someFunction(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

## Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported functions
- Update configuration examples if needed
- Include code examples in documentation

## Testing Your Changes

Before submitting a pull request:

1. Run all tests:
```bash
make test
```

2. Check test coverage:
```bash
make test-coverage
```

3. Format your code:
```bash
make fmt
```

4. Run vet:
```bash
make vet
```

5. Build the application:
```bash
make build
```

6. Test manually if applicable:
```bash
make run
```

## Review Process

1. All pull requests require review
2. Address review comments promptly
3. Keep discussions focused and respectful
4. Be open to feedback and suggestions

## Release Process

Releases are managed by maintainers:

1. Version bumps follow [Semantic Versioning](https://semver.org/)
2. Changelog is updated for each release
3. Git tags are created for releases
4. GitHub releases include binaries

## Questions?

If you have questions about contributing:

1. Check existing issues and discussions
2. Read the project documentation
3. Create a new issue with the `question` label

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to OSS Insight Trending Notifier!
