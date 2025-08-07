# Suggested Commands

## Development Commands

### Building
```bash
make build          # Build the binary as 'f'
make all           # Clean and build
make clean         # Remove build artifacts
```

### Testing
```bash
make test          # Run all tests with verbose output
go test -v ./...   # Alternative test command
```

### Development Setup
```bash
make dev-setup     # Resolve and verify dependencies
go mod tidy        # Clean up go.mod
go mod verify      # Verify dependencies
```

### Version Information
```bash
make version       # Show version and build time
```

### Installation
```bash
make install       # Install binary to /usr/local/bin/
make uninstall     # Remove installed binary
```

### Running the Tool
```bash
./f               # Run built binary (shows all commands)
./f test          # Filter commands containing "test"
./f "run dev"     # Filter commands containing "run dev"
```

## Git Commands
```bash
git status        # Check repository status
git add .         # Stage all changes
git commit -m ""  # Commit with message
git push          # Push to remote
```

## System Commands (Darwin/macOS)
```bash
ls -la            # List all files including hidden
find . -name "*"  # Find files
grep -r "text"    # Search in files recursively
which command     # Find command location
```

## Go Specific Commands
```bash
go version        # Check Go version
go run main.go    # Run without building
go build          # Build current package
go get -u         # Update dependencies
```

## External Dependencies Check
```bash
which fzf         # Verify fzf is installed
brew install fzf  # Install fzf if missing (macOS)
```