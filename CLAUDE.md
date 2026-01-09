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

## Common Commands

### Building & Running
```bash
# Build the project (to be added when implemented)
go build -o reactor ./cmd/reactor

# Run the application (to be added when implemented)
go run ./cmd/reactor

# Run with verbose logging (to be added when implemented)
go run ./cmd/reactor -v
```

### Testing
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

### Development
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
