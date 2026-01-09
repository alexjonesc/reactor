# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Reactor - MIDI Processing System

## Role & Expertise
You are a senior GoLang engineer with specialized expertise in MIDI music and MIDI instruments. You understand MIDI protocols, note structures, timing, and instrument communication.

## Project Overview
This is a Go-based MIDI processing system that receives MIDI inputs from various devices, processes them through configurable algorithms, and sends MIDI outputs to available ports.

## Development Phases

### Phase 1: MIDI Input Processing (Current)
- Implement MIDI input handling from various devices
- Process inputs through algorithmic transformations
- Output results in text format for debugging and validation

### Phase 2: Clock System
- Implement a controllable click/metronome system with three modes:
  - Simple pulse-based timing
  - MIDI clock input synchronization
  - Manual stepping via specific MIDI note inputs

### Phase 3: MIDI Output
- Implement full MIDI output to available ports
- Send processed MIDI sequences to external devices

## Technical Principles

### Architecture
- **Composability**: Algorithms should be chainable in different configurations
- **Configurability**: Settings-driven algorithm parameters
- **Simplicity**: Prioritize clean, readable code over clever solutions
- **Testability**: Write code that's easy to unit test and validate
- **Reusability**: Design components that can be used in multiple contexts

### Code Quality Standards
- Write idiomatic Go code
- Keep functions small and focused
- Use clear, descriptive names
- Minimize dependencies between components
- Document MIDI-specific logic and timing considerations

## Development Approach
Start with the simplest working implementation, then iterate. Build each phase completely before moving to the next. Focus on getting MIDI input processing solid before tackling timing and output.

## Docker Development Environment

This project uses Docker with Air (hot reload) and Delve (debugging) for development.

### Prerequisites
- Docker Desktop 4.35+ (for USB/IP MIDI device support on macOS)
- Docker Compose V2
- Make (optional, for convenience commands)
- Go 1.25+ (for native development outside Docker)

### Quick Start
```bash
# First time setup
make env          # Create .env from .env.example
make build-dev    # Build development Docker image

# Start development environment with hot reload
make run-dev

# In another terminal, attach VS Code debugger
# Run > Start Debugging > "Connect to Docker Delve" (localhost:2345)
```

### Common Docker Commands

**Development:**
```bash
make run-dev      # Start with Air hot reload and Delve debugger
make debug        # Start in background (daemon mode)
make stop         # Stop all containers
make logs         # Follow container logs
make shell        # Open shell in running container
make clean        # Remove containers, volumes, and tmp files
```

**Building:**
```bash
make build-dev    # Build development image (with Air + Delve)
make build-prod   # Build production image (optimized, ~15MB)
```

**Native Go Commands (when not using Docker):**
```bash
go build -o reactor .
go run .
go test ./...
```

### Development Workflow

**Hot Reload Development:**
1. Start the environment: `make run-dev`
2. Edit any `.go` file in your IDE
3. Air automatically detects changes, rebuilds with debug flags, and restarts
4. Delve remains attached and ready for debugging

**Debugging:**
1. Start in background: `make debug`
2. In VS Code: Run > Start Debugging > "Connect to Docker Delve"
3. Set breakpoints in your code
4. Trigger MIDI events or application logic
5. Debugger pauses at breakpoints

**Testing:**
```bash
make test         # Run all tests in container
docker-compose run --rm reactor-dev go test ./clock -v  # Run specific package
```

### MIDI Device Setup (macOS)

**Challenge:** Docker on macOS uses a Linux VM, making MIDI device passthrough complex.

**Option 1: Docker Desktop USB Sharing (Recommended for Physical Devices)**
1. Open Docker Desktop > Settings > Resources > USB Sharing
2. Enable USB sharing for your MIDI devices (e.g., SE25)
3. Restart containers if needed
4. Update `.env` file with correct device names

**Option 2: Native Development (Recommended for Virtual MIDI)**
- IAC Driver (virtual MIDI bus) may not be accessible via USB/IP
- Run application natively on macOS for development: `go run .`
- Use Docker for CI/CD and production builds only
- Environment variables still work: `export MIDI_INPUT_DEVICE=SE25 && go run .`

**Option 3: Network MIDI (Advanced)**
- Use RTP-MIDI or ipMIDI to bridge host and container
- More complex setup, but works for both physical and virtual devices

### Environment Configuration

Edit `.env` to match your MIDI setup:
```bash
MIDI_INPUT_DEVICE=SE25                # Your MIDI input device name
MIDI_OUTPUT_DEVICE=IAC Driver Bus 1   # Your MIDI output device name
DEBUG_MODE=true
LOG_LEVEL=info
```

To list available MIDI devices on your system, you'll need to run code that queries MIDI ports.

### Troubleshooting

**Delve won't start:**
- Ensure `security_opt: seccomp:unconfined` is in docker-compose.yml ✓ (already configured)
- Check port 2345 isn't in use: `lsof -i :2345`
- Check container logs: `make logs`

**Hot reload not working:**
- Ensure source code is mounted: Check `docker-compose.yml` volumes ✓ (already configured)
- Verify Air is running: `make logs` should show Air watching files
- Check `.air.toml` includes correct directories

**MIDI devices not found:**
- Verify Docker Desktop USB sharing is enabled (Settings > Resources > USB Sharing)
- Check device names match your system exactly (case-sensitive)
- For IAC Driver (virtual), consider running natively instead of Docker
- Test with: `make shell` then check MIDI device detection in code

**Build errors:**
- Clean and rebuild: `make clean && make build-dev`
- Check Go module cache: Docker uses a named volume for `/go/pkg/mod`
- Verify `go.mod` and `go.sum` are present

### Go Development Commands (Native)

When running natively without Docker:

**Testing:**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestName ./path/to/package

# Run tests with verbose output
go test -v ./...
```

**Development Tools:**
```bash
# Format code
go fmt ./...

# Run linter (if configured)
golangci-lint run

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Key Implementation Notes

### MIDI Processing Pipeline
The system follows a pipeline architecture:
1. **Input Layer**: Receives MIDI events from devices
2. **Processing Layer**: Applies configurable algorithms to transform MIDI data
3. **Output Layer**: Sends results to MIDI ports or text output (phase-dependent)

### Algorithm Composability
Algorithms should implement a common interface to enable chaining:
- Each algorithm receives MIDI events and produces transformed events
- Algorithms can be configured via settings/configuration files
- Chain execution order matters for musical results

### MIDI-Specific Considerations
- **Timing**: MIDI timing is critical; avoid blocking operations in the event loop
- **Note On/Off Pairing**: Ensure note-on events have corresponding note-off events
- **Velocity Sensitivity**: Preserve or transform velocity data appropriately
- **Channel Awareness**: Track which MIDI channel events belong to
