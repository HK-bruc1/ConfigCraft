# Contributing to ConfigCraft

Thank you for your interest in contributing to ConfigCraft! We welcome contributions from the community and are pleased to have them.

## Quick Start

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/configcraft.git`
3. **Create** a new branch: `git checkout -b feature/your-feature-name`
4. **Make** your changes
5. **Test** your changes: `go test ./...`
6. **Commit** your changes: `git commit -m 'Add some feature'`
7. **Push** to your fork: `git push origin feature/your-feature-name`
8. **Submit** a Pull Request

## Development Guidelines

### Code Style
- Follow Go conventions and use `gofmt`
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused

### Version Updates
- Only modify `internal/version/version.go` for version changes
- Version format: `MAJOR.MINOR.PATCH`

### UI Components
- Follow existing patterns in `internal/ui/components/`
- Ensure cross-platform compatibility
- Test with different window sizes

### Documentation
- Update relevant documentation for any API changes
- Add examples for new features
- Keep CHANGELOG.md updated with technical details

## Reporting Issues

Please use the [GitHub issue tracker](https://github.com/ConfigCraft/configcraft/issues) to report bugs or request features.

### Bug Reports
Include:
- Operating system and version
- Go version (if building from source)
- Steps to reproduce the issue
- Expected vs actual behavior
- Error logs if available

### Feature Requests
Include:
- Use case and motivation
- Detailed description of the feature
- Examples of how it would be used

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## License

By contributing to ConfigCraft, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing! 🎉