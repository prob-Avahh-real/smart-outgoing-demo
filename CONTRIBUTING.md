# Contributing Guidelines

## Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd smart-outgoing-demo
   ```

2. **Install dependencies**
   ```bash
   # Go dependencies
   go mod download
   
   # Frontend dependencies (if any)
   npm install
   ```

3. **Environment setup**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## Code Standards

### Go Best Practices

- **Package structure**: Follow standard Go project layout
- **Error handling**: Use explicit error types and proper error wrapping
- **Logging**: Use structured logging with the provided logger package
- **Testing**: Write unit tests for all public functions
- **Documentation**: Add godoc comments for all exported functions
- **Linting**: Run `golangci-lint run` before committing

### JavaScript Best Practices

- **ES6+ syntax**: Use modern JavaScript features
- **Modules**: Use ES6 modules and proper imports/exports
- **Error handling**: Use try/catch blocks and proper error propagation
- **Code organization**: Separate concerns into different classes/modules
- **Performance**: Optimize for performance and memory usage

### General Guidelines

- **Commit messages**: Use conventional commit format
- **Branch naming**: Use descriptive branch names (feature/xxx, bugfix/xxx)
- **Code reviews**: All changes must be reviewed before merging
- **Testing**: Ensure all tests pass before submitting PR

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/store

# Run tests with verbose output
go test -v ./...
```

### Test Coverage

Maintain at least 80% test coverage for all packages.

## Code Quality

### Linting

```bash
# Run Go linter
golangci-lint run

# Format code
go fmt ./...
```

### Static Analysis

```bash
# Run security checks
gosec ./...

# Run vulnerability checks
govulncheck ./...
```

## Pull Request Process

1. Create a new branch from `main`
2. Make your changes
3. Add tests for new functionality
4. Ensure all tests pass
5. Run linting and fix any issues
6. Submit a pull request with:
   - Clear description of changes
   - Testing instructions
   - Any breaking changes

## Architecture Guidelines

### Backend

- **Layers**: Handler -> Service -> Repository pattern
- **Dependencies**: Use dependency injection
- **Configuration**: Environment-based configuration
- **Error handling**: Consistent error types across layers

### Frontend

- **Components**: Modular component-based architecture
- **State management**: Centralized state management
- **API communication**: Dedicated service layer
- **Error handling**: User-friendly error messages

## Security

- **Input validation**: Validate all user inputs
- **Authentication**: Proper token-based authentication
- **Authorization**: Role-based access control
- **Data sanitization**: Sanitize all data before processing

## Performance

- **Caching**: Implement appropriate caching strategies
- **Database**: Optimize database queries
- **Memory**: Monitor and optimize memory usage
- **Concurrency**: Use goroutines for concurrent operations

## Documentation

- **API docs**: Keep API documentation up to date
- **Code comments**: Add meaningful comments
- **README**: Update README with new features
- **Changelog**: Maintain changelog for releases
