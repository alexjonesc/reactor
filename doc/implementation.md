# Docker + Delve + Air Setup for Reactor MIDI Application

## Overview
Set up a complete Docker development environment on the `poc-2` branch with:
- Multi-stage Dockerfile (dev stage with Air + Delve, production stage optimized)
- Hot reloading via Air
- Remote debugging via Delve on port 2345
- MIDI device access (SE25 input, IAC Driver output)
- Environment-based configuration for MIDI devices and debug settings

## Implementation Steps

### 1. Create Docker Infrastructure Files

**File: `Dockerfile`**
- Base stage: Go 1.24.5-alpine, copy go.mod/go.sum, run `go mod download`
- Development stage: Install Air and Delve, copy source, expose port 2345, CMD runs Air
- Builder stage: Build optimized binary with `CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s"`
- Production stage: Minimal alpine image with compiled binary only

**File: `docker-compose.yml`**
- Service: `reactor-dev` targeting development stage
- Security options: `seccomp:unconfined` and `SYS_PTRACE` capability (required for Delve)
- Volumes: Mount source code (`.:/app`) and Go modules cache
- Ports: Expose 2345 for Delve debugger
- Environment: Load from `.env` file
- Note: MIDI device access on macOS may require Docker Desktop 4.35+ USB/IP feature or host network mode

**File: `.air.docker.toml`**
- Build command: `go build -gcflags='all=-N -l' -o ./tmp/main .`
- Run command: `dlv exec --headless --listen=:2345 --api-version=2 --accept-multiclient --continue ./tmp/main`
- Watch: `*.go` files in all directories, exclude `tmp/`, `.git/`, `vendor/`
- Delay: 1000ms rebuild delay
- Result: Air watches files → rebuilds with debug flags → launches via Delve

**File: `.env.example`**
```
MIDI_INPUT_DEVICE=SE25
MIDI_OUTPUT_DEVICE=IAC Driver Bus 1
DEBUG_MODE=true
LOG_LEVEL=info
```

**File: `.dockerignore`**
- Exclude: `.git/`, `*.md` (except README), IDE files, build artifacts, `.env` (keep `.env.example`)
- Purpose: Minimize build context, speed up builds

### 2. Create Developer Experience Files

**File: `Makefile`**
Common commands:
- `make build-dev` - Build dev image
- `make run-dev` - Start with hot reload
- `make debug` - Start in background for debugging
- `make stop` - Stop containers
- `make clean` - Remove containers and volumes
- `make logs` - Follow logs
- `make shell` - Open shell in container
- `make test` - Run tests
- `make build-prod` - Build production image

**File: `.vscode/launch.json`**
VS Code debugging configuration:
- Type: Go Remote
- Host: localhost, Port: 2345
- Remote path mapping: `/app` → local repo path

### 3. Create Initial Go Project Structure

Starting from scratch on poc-2 branch. When creating new Go code:

**File: `go.mod`**
```
module reactor
go 1.24.5
require gitlab.com/gomidi/midi/v2 v2.3.16
```

**File: `main.go`** (to be created when implementing application)
Use environment variables from the start:
```go
import "os"

inputDevice := os.Getenv("MIDI_INPUT_DEVICE")
if inputDevice == "" {
    inputDevice = "SE25" // default
}

outputDevice := os.Getenv("MIDI_OUTPUT_DEVICE")
if outputDevice == "" {
    outputDevice = "IAC Driver Bus 1" // default
}

in, err := midi.FindInPort(inputDevice)
out, err := midi.FindOutPort(outputDevice)
```

### 4. Update CLAUDE.md

Add Docker commands section with:
- Quick start instructions (`cp .env.example .env`, `make run-dev`)
- Development workflow (hot reload, debugging)
- MIDI device setup for macOS (Docker Desktop USB sharing)
- Testing and troubleshooting tips
- Common commands reference

### 5. Testing & Validation

**Test hot reload:**
1. `make run-dev`
2. Edit any .go file
3. Verify Air detects change, rebuilds, restarts

**Test debugging:**
1. `make debug` (starts in background)
2. Connect IDE debugger to localhost:2345
3. Set breakpoint in main.go
4. Verify debugger hits breakpoint

**Test MIDI (if devices available):**
1. Enable Docker Desktop USB sharing (Settings → Resources → USB Sharing)
2. Share SE25 and IAC Driver devices
3. Verify application can detect and use devices

## Critical Files to Create/Modify

New files (poc-2 branch):
1. `Dockerfile` - Multi-stage build definition
2. `docker-compose.yml` - Dev environment orchestration
3. `.air.docker.toml` - Hot reload configuration
4. `.env.example` - Environment variable template
5. `.dockerignore` - Build optimization
6. `Makefile` - Developer commands
7. `.vscode/launch.json` - IDE debugging config

Modified files:
8. `CLAUDE.md` - Add Docker commands and workflow documentation

Go application files (to be created as needed):
- `go.mod`, `go.sum` - Go module definition
- `main.go` - Application entry point with environment variable support
- Application code (clock/, sequencer/, etc.) as development progresses

## MIDI Device Access Considerations (macOS)

**Challenge:** Docker on macOS uses a Linux VM, making device passthrough complex.

**Solutions (in order of preference):**
1. **Docker Desktop 4.35+ USB/IP**: Use built-in USB sharing feature for physical devices
2. **Host Network Mode**: Not available on Docker Desktop for Mac
3. **Native Development**: Run application natively on macOS, use Docker for CI/CD only
4. **Network MIDI**: Bridge host and container via RTP-MIDI (complex setup)

**Recommendation:** Test USB/IP first. If IAC Driver (virtual device) isn't accessible, consider native development workflow while keeping Docker for production builds.

## Implementation Order

1. Create `.dockerignore` (fast, no dependencies)
2. Create `Dockerfile` (core infrastructure)
3. Create `docker-compose.yml` (depends on Dockerfile)
4. Create `.air.docker.toml` (Air config)
5. Create `.env.example` (configuration template)
6. Create `Makefile` (convenience layer)
7. Create `go.mod` (Go module initialization)
8. Create `.vscode/launch.json` (debugging support)
9. Update `CLAUDE.md` documentation
10. Test infrastructure: `make build-dev` (should build successfully)

**Note:** The Go application code (main.go, packages, etc.) will be developed separately. The Docker/Delve/Air infrastructure will be ready to support development from the start.

## Success Criteria

- ✅ `make run-dev` starts Air watching Go files
- ✅ Code changes trigger automatic rebuild and restart
- ✅ Debugger accessible on localhost:2345 from IDE
- ✅ Environment variables control MIDI device names
- ✅ Production build creates optimized binary
- ✅ Documentation covers all common workflows
