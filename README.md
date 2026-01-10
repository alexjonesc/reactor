# Reactor - MIDI Processing System

A Go-based MIDI processing system that receives MIDI inputs from various devices, processes them through configurable algorithms, and sends MIDI outputs to available ports.

## Quick Start

### Choose Your Development Environment:

**Native (Recommended for MIDI):**
```bash
make run-native     # Full MIDI access, hot reload, debugging
```

**Docker:**
```bash
make run-dev        # Containerized, consistent environment
```

## Features

- ✅ **MIDI Input Handling** - Comprehensive event processing
  - Note On/Off with note names (C4, A#3, etc.)
  - Control Change (CC) messages
  - Pitch Bend, Program Change, Aftertouch
  - System Exclusive (SysEx) messages

- ✅ **Hot Reload** - Air watches for code changes
- ✅ **Remote Debugging** - Delve integration with VS Code
- ✅ **Environment Configuration** - Device names via `.env`
- ✅ **Graceful Fallback** - Runs without MIDI devices

## Development Environments

### Native Development (for MIDI Testing)

**Best for:** Working with MIDI hardware, IAC Driver, testing algorithms

**Features:**
- ✅ Full MIDI device access (physical + virtual)
- ✅ Hot reload with Air
- ✅ Debugging with Delve (port 2346)
- ✅ Lower latency

**Setup:**
```bash
# Install Air (if needed)
go install github.com/air-verse/air@latest

# Start development
make run-native
```

**Documentation:** See `NATIVE_DEVELOPMENT.md`

### Docker Development

**Best for:** General development, CI/CD, consistent environments

**Features:**
- ✅ Hot reload with Air
- ✅ Debugging with Delve (port 2345)
- ⚠️ Limited MIDI access on macOS

**Setup:**
```bash
make env         # Create .env file
make build-dev   # Build Docker image
make run-dev     # Start development
```

**Documentation:** See `CLAUDE.md` for complete Docker setup

## MIDI Device Support

### Tested Devices:
- ✅ SE25 MIDI keyboard (physical USB device)
- ✅ IAC Driver Bus 1 (virtual MIDI bus - macOS)

### Configuration:

Edit `.env` file:
```env
MIDI_INPUT_DEVICE=SE25
MIDI_OUTPUT_DEVICE=IAC Driver Bus 1
DEBUG_MODE=true
LOG_LEVEL=info
```

## Project Structure

```
reactor/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── .env.example               # Environment variables template
├── .air.docker.toml           # Docker Air configuration
├── .air.native.toml           # Native Air configuration
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yml         # Docker orchestration
├── Makefile                   # Development commands
├── CLAUDE.md                  # Complete development guide
├── NATIVE_DEVELOPMENT.md      # Native setup guide
├── MIDI_GUIDE.md              # MIDI implementation details
└── DEBUG_GUIDE.md             # Debugging instructions
```

## Common Commands

### Native Development:
```bash
make run-native       # Start with hot reload
make debug-native     # Start in background for debugging
make stop-native      # Stop background process
make test-native      # Run tests
```

### Docker Development:
```bash
make run-dev          # Start with hot reload
make debug            # Start in background
make stop             # Stop containers
make logs             # View logs
make shell            # Open shell in container
make test             # Run tests in container
```

### Both:
```bash
make env              # Create .env from template
make help             # Show all available commands
```

## Example Output

```
🎵 Reactor MIDI Processing System Starting...
===================================================
Configuration:
  MIDI Input:  SE25
  MIDI Output: IAC Driver Bus 1
  Debug Mode:  true
  Log Level:   info
===================================================

📋 Available MIDI Ports:
───────────────────────────────────────────────────
Input Ports (3):
  1. IAC Driver Bus 1
  2. SE25 MIDI1
  3. SE25 MIDI2
───────────────────────────────────────────────────

🎹 Connecting to MIDI input: SE25
✅ Connected to: SE25 MIDI1

🎵 Listening for MIDI input... (Press Ctrl+C to stop)

[14:23:15.234] 🎹 Note ON  | Ch: 0 | Note:  60 (C4) | Vel: 100
[14:23:15.456] 🎹 Note OFF | Ch: 0 | Note:  60 (C4)
[14:23:16.123] 🎹 Note ON  | Ch: 0 | Note:  64 (E4) | Vel:  85
[14:23:20.456] 🎛️  CC       | Ch: 0 | CC:   1 | Val:  64
```

## Development Phases

### ✅ Phase 1: MIDI Input Processing (Current)
- Implement MIDI input handling from various devices
- Process inputs through algorithmic transformations
- Output results in text format for debugging and validation

### 🔄 Phase 2: Clock System (Next)
- Implement a controllable click/metronome system with three modes:
  - Simple pulse-based timing
  - MIDI clock input synchronization
  - Manual stepping via specific MIDI note inputs

### 📋 Phase 3: MIDI Output (Planned)
- Implement full MIDI output to available ports
- Send processed MIDI sequences to external devices

## Technical Stack

- **Language:** Go 1.25
- **MIDI Library:** gitlab.com/gomidi/midi/v2
- **Hot Reload:** Air v1.63.6
- **Debugger:** Delve v1.26.0
- **Containerization:** Docker + Docker Compose

## Requirements

- Go 1.25 or later
- macOS (for native MIDI development)
- Docker Desktop 4.35+ (for Docker development)
- MIDI device (keyboard, controller, or virtual device)

## Documentation

- **`CLAUDE.md`** - Complete development guide and project overview
- **`NATIVE_DEVELOPMENT.md`** - Native setup with Air and Delve
- **`MIDI_GUIDE.md`** - MIDI implementation and testing
- **`DEBUG_GUIDE.md`** - VS Code debugging setup
- **`implementation.md`** - Docker infrastructure plan

## Troubleshooting

### MIDI device not found:
```bash
# Run to see available devices
make run-native

# Update .env with exact device name from list
```

### Port already in use:
```bash
# Native (port 2346)
make stop-native

# Docker (port 2345)
make stop
```

### Hot reload not working:
```bash
# Check Air is running
ps aux | grep air

# Restart
make stop-native && make run-native
```

See individual guides for detailed troubleshooting.

## Contributing

This is a personal project, but suggestions and feedback are welcome!

## License

[Add your license here]

---

**Start developing:**

```bash
# For MIDI testing
make run-native

# For general development
make run-dev
```

🎵 Happy coding!
