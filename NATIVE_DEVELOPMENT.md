# Native Development Guide

This guide explains how to develop Reactor natively on macOS with full MIDI device access, hot reload (Air), and debugging (Delve).

## Status: ✅ READY

Native development is fully configured and tested!

## Why Native Development?

### Benefits:
- ✅ **Full MIDI Access** - All physical and virtual devices (SE25, IAC Driver, etc.)
- ✅ **Lower Latency** - Direct hardware access, no Docker overhead
- ✅ **Hot Reload** - Air watches for changes and rebuilds automatically
- ✅ **Remote Debugging** - Delve on port 2346, VS Code integration
- ✅ **Simpler Setup** - No container complexity for MIDI testing

### When to Use:
- Testing with real MIDI devices
- Developing MIDI processing algorithms
- Testing IAC Driver (virtual MIDI bus)
- Low-latency requirements
- Quick iteration on MIDI features

## Quick Start

### 1. Start Development Server

**Option A: Foreground (see output directly)**
```bash
make run-native
```

**Option B: Background (for debugging)**
```bash
make debug-native
```

### 2. Connect Debugger (Optional)

If using `debug-native`:

1. Open VS Code
2. Go to **Run and Debug** (Cmd+Shift+D)
3. Select **"Connect to Native Delve (Air)"**
4. Click the green play button
5. Set breakpoints and debug!

### 3. Stop Server

```bash
make stop-native
```

Or press **Ctrl+C** if running in foreground.

## Available Commands

```bash
make run-native      # Start with hot reload (foreground)
make debug-native    # Start in background for debugging
make stop-native     # Stop background server
make test-native     # Run tests natively
```

## Configuration

### Environment Variables

Same `.env` file as Docker:

```env
MIDI_INPUT_DEVICE=SE25
MIDI_OUTPUT_DEVICE=IAC Driver Bus 1
DEBUG_MODE=true
LOG_LEVEL=info
```

### Air Configuration

File: `.air.native.toml`

Key settings:
- **Debugger Port**: 2346 (avoids conflict with Docker's 2345)
- **Build Command**: `go build -gcflags='all=-N -l'` (debug flags)
- **Watch**: All `.go` files
- **Excludes**: `tmp/`, `.git/`, `vendor/`

## Debugging

### VS Code Setup

Two debug configurations available:

1. **"Connect to Docker Delve"** - Port 2345 (Docker)
2. **"Connect to Native Delve (Air)"** - Port 2346 (Native)

### Debug Workflow

1. **Start native server:**
   ```bash
   make debug-native
   ```

2. **Connect VS Code debugger:**
   - Run > Start Debugging
   - Select "Connect to Native Delve (Air)"
   - Wait for connection (~2 seconds)

3. **Set breakpoints:**
   - Open `main.go`
   - Click in the gutter to set breakpoints
   - Good locations:
     - Line 100: `handleMIDIMessage()` - catches all MIDI events
     - Line 103: Note On handler
     - Line 40: Device connection

4. **Test with MIDI device:**
   - Play notes on your MIDI keyboard
   - Debugger pauses at breakpoints
   - Inspect variables, step through code

5. **Hot reload while debugging:**
   - Edit code while debugger is connected
   - Air detects changes and rebuilds
   - Debugger automatically reconnects
   - Breakpoints remain active!

### Debug Console Commands

When connected, use the Debug Console:

```
dlv> help            # Show available commands
dlv> goroutines      # List all goroutines
dlv> stack           # Show call stack
dlv> print msg       # Print MIDI message variable
dlv> locals          # Show local variables
dlv> continue        # Resume execution
```

## MIDI Device Testing

### Expected Output

```bash
make run-native
```

You should see:

```
  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ v1.63.6, built with Go go1.25.5

[09:49:39] watching .
[09:49:39] !exclude tmp
[09:49:39] building...
[09:49:44] running...

API server listening at: [::]:2346
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

Output Ports (2):
  1. IAC Driver Bus 1
  2. SE25 MIDI1
───────────────────────────────────────────────────

🎹 Connecting to MIDI input: SE25
✅ Connected to: SE25 MIDI1

🎵 Listening for MIDI input... (Press Ctrl+C to stop)
```

### Testing MIDI Events

Play notes on your MIDI keyboard:

```
[14:23:15.234] 🎹 Note ON  | Ch: 0 | Note:  60 (C4) | Vel: 100
[14:23:15.456] 🎹 Note OFF | Ch: 0 | Note:  60 (C4)
[14:23:16.123] 🎹 Note ON  | Ch: 0 | Note:  64 (E4) | Vel:  85
```

Move knobs/faders:

```
[14:23:20.456] 🎛️  CC       | Ch: 0 | CC:   1 | Val:  64
```

Pitch bend:

```
[14:23:21.123] 🎚️  Bend     | Ch: 0 | Value:   2048
```

## Comparing Docker vs Native

| Feature | Docker | Native |
|---------|--------|--------|
| Hot Reload | ✅ Air | ✅ Air |
| Debugging | ✅ Port 2345 | ✅ Port 2346 |
| Physical MIDI | ⚠️ USB Sharing | ✅ Direct Access |
| Virtual MIDI | ❌ Not Available | ✅ Full Access |
| IAC Driver | ❌ Not Available | ✅ Works Perfectly |
| Setup | Complex | Simple |
| Latency | Higher | Lower |
| Best For | General Dev | MIDI Testing |

## Workflow Recommendations

### Development Phase Approach:

**Writing Code (Use Docker):**
```bash
make run-dev        # Docker with hot reload
```
- Good for: Writing non-MIDI code, tests, refactoring
- Benefits: Consistent environment, CI/CD ready

**Testing MIDI (Use Native):**
```bash
make run-native     # Native with full MIDI access
```
- Good for: Testing with real hardware, MIDI algorithms
- Benefits: All devices accessible, lower latency

**Debugging MIDI (Use Native):**
```bash
make debug-native   # Native with Delve + MIDI access
```
- Good for: Debugging MIDI event handling, timing issues
- Benefits: Full MIDI + debugger + hot reload

## File Structure

```
reactor/
├── .air.toml               # Docker Air config (port 2345)
├── .air.native.toml        # Native Air config (port 2346)
├── .vscode/
│   └── launch.json         # Both Docker and Native configs
├── tmp/
│   ├── main                # Compiled binary
│   ├── air.log             # Air logs (background mode)
│   └── air.pid             # Process ID (background mode)
└── Makefile                # Docker and Native commands
```

## Troubleshooting

### Port 2346 Already in Use

**Cause:** Previous Air process still running

**Solution:**
```bash
make stop-native
# Or manually
pkill -f "dlv.*2346"
lsof -ti:2346 | xargs kill -9
```

### Air Command Not Found

**Cause:** Air not installed or not in PATH

**Solution:**
```bash
go install github.com/air-verse/air@latest

# Verify installation
air -v

# Add to PATH if needed
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Delve Command Not Found

**Cause:** Delve not installed

**Solution:**
```bash
go install github.com/go-delve/delve/cmd/dlv@latest

# Verify installation
dlv version
```

### MIDI Device Not Found

**Cause:** Device name mismatch or device not connected

**Solution:**
1. Run application to see available devices:
   ```bash
   make run-native
   ```

2. Check "Available MIDI Ports" section

3. Update `.env` with exact device name:
   ```env
   MIDI_INPUT_DEVICE=SE25 MIDI1
   ```

### Hot Reload Not Working

**Cause:** Air not detecting changes

**Solution:**
1. Check Air is watching: Look for `[watching .]` in output
2. Verify `.air.native.toml` exists
3. Check file is not in excluded directories
4. Try saving the file again

### Debugger Won't Connect

**Cause:** Delve not running or wrong port

**Solution:**
1. Verify Air is running: `ps aux | grep air`
2. Check Delve is listening: `lsof -i :2346`
3. Check VS Code config uses port 2346
4. Try reconnecting debugger

### Build Errors with rtmidi

**Cause:** CGO compilation warnings (safe to ignore)

**These warnings are normal:**
```
RtMidi.cpp:1685:15: warning: variable length arrays in C++ are a Clang extension
```

These don't affect functionality.

## Advanced: Running Without Air

For production-like testing:

```bash
# Build with debug flags
go build -gcflags='all=-N -l' -o reactor .

# Run with Delve manually
dlv exec --headless --listen=:2346 --api-version=2 ./reactor

# Or just run directly
./reactor
```

## Performance Tips

### Faster Rebuilds

Air is already optimized, but you can:

1. **Exclude more directories** in `.air.native.toml`:
   ```toml
   exclude_dir = ["tmp", "vendor", ".git", ".idea", ".vscode", "docs"]
   ```

2. **Disable unnecessary watches**:
   ```toml
   exclude_regex = ["_test.go", ".*_test.go"]
   ```

3. **Reduce rebuild delay**:
   ```toml
   delay = 500  # milliseconds
   ```

### Lower MIDI Latency

For time-critical MIDI processing:

1. Run without debugger:
   ```bash
   go run .
   ```

2. Build optimized binary:
   ```bash
   go build -o reactor .
   ./reactor
   ```

3. Use Release build flags:
   ```bash
   go build -ldflags="-w -s" -o reactor .
   ./reactor
   ```

## Next Steps

Now that native development is set up:

1. **Connect your MIDI device** (SE25, keyboard, controller)
2. **Start development**: `make run-native`
3. **Play some notes** - watch events in real-time
4. **Set breakpoints** in `handleMIDIMessage()`
5. **Edit code** - watch Air rebuild automatically
6. **Build MIDI algorithms** with full hardware access!

---

## Summary

✅ **Installed**: Air v1.63.6, Delve v1.26.0
✅ **Configured**: `.air.native.toml`, VS Code launch config
✅ **Tested**: MIDI devices detected (SE25, IAC Driver)
✅ **Ready**: Hot reload + Debugging + Full MIDI access

**Start developing:**
```bash
make run-native
```

🎵 Happy MIDI coding!
