# Task Completion Workflow

## When a Task is Completed

### 1. Code Quality Checks
```bash
go fmt ./...       # Format all Go code
go vet ./...       # Run Go vet for potential issues
```

### 2. Testing
```bash
make test          # Run all tests
# or
go test -v ./...   # Alternative test command
```

### 3. Build Verification
```bash
make build         # Ensure code builds successfully
# or
make all          # Clean build from scratch
```

### 4. Dependencies Check
```bash
go mod tidy        # Clean up go.mod/go.sum
go mod verify      # Verify dependency checksums
```

### 5. Static Analysis (Optional)
```bash
go vet ./...       # Built-in static analysis
```

### 6. Manual Testing (if applicable)
```bash
./f               # Test the built binary
./f test          # Test filtering functionality
```

## Pre-commit Checklist
- [ ] Code formatted with `go fmt`
- [ ] All tests pass with `make test`
- [ ] Code builds successfully with `make build`
- [ ] Dependencies are clean with `go mod tidy`
- [ ] Manual testing completed (if UI changes)
- [ ] Comments added in Japanese for complex logic
- [ ] Error messages are in English

## Common Issues to Check
- External dependency availability (fzf)
- Cross-platform compatibility considerations
- Plugin interface compliance
- Configuration file handling
- Error message clarity and language

## Build Artifacts
- Binary name: `f`
- Clean up with: `make clean`
- Installation location: `/usr/local/bin/f`