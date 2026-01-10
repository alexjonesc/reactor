# MIDI Input Handling Guide

This guide explains the MIDI input implementation and how to test it.

## Implementation Status: ✅ COMPLETE

Basic MIDI input handling has been implemented with the following features:

### Features Implemented

1. **MIDI Device Discovery**
   - Lists all available MIDI input and output ports on startup
   - Displays port names and connection status

2. **Environment-Based Configuration**
   - `MIDI_INPUT_DEVICE` - Configure which MIDI input to use (default: "SE25")
   - `MIDI_OUTPUT_DEVICE` - Configure MIDI output device (default: "IAC Driver Bus 1")

3. **MIDI Event Handling**
   - ✅ Note On/Off events with note names (e.g., "C4", "A#3")
   - ✅ Control Change (CC) messages
   - ✅ Pitch Bend
   - ✅ Program Change
   - ✅ Channel Aftertouch
   - ✅ Polyphonic Aftertouch
   - ✅ System Exclusive (SysEx) messages

4. **Graceful Fallback**
   - Runs without MIDI devices if none are available
   - Clear error messages when configured device not found
   - Continues running in fallback mode with status updates

5. **Real-Time Event Display**
   - Timestamps for all MIDI events (millisecond precision)
   - Color-coded event types with emojis
   - Channel, note, velocity information
   - Human-readable note names

## Testing in Docker (Expected Behavior)

### What You'll See:

```bash
make debug
# Wait for startup, then check logs:
make logs
```

**Expected Output:**
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
Input Ports (0):
  (none)

Output Ports (0):
  (none)
───────────────────────────────────────────────────

🎹 Connecting to MIDI input: SE25
❌ Error: Could not find MIDI input device 'SE25'

ℹ️  Running without MIDI input. Check available ports above.

⚠️  Running in fallback mode (no MIDI input)
```

**This is expected!** MIDI devices are not accessible from Docker containers on macOS.

## Testing with Real MIDI Devices (Native Development)

To test with actual MIDI hardware, run natively on macOS:

### Setup:

1. **Connect your MIDI device** (e.g., MIDI keyboard, controller)

2. **Verify device is connected:**
   ```bash
   # macOS has built-in MIDI support
   # Your device should be automatically recognized
   ```

3. **Configure environment:**
   ```bash
   # Edit .env file to match your MIDI device name
   # Or set environment variables:
   export MIDI_INPUT_DEVICE="Your MIDI Device Name"
   export MIDI_OUTPUT_DEVICE="IAC Driver Bus 1"
   ```

4. **Run natively:**
   ```bash
   go run .
   ```

### Expected Output (with MIDI device):

```
🎵 Reactor MIDI Processing System Starting...
===================================================
Configuration:
  MIDI Input:  Keystation 49
  MIDI Output: IAC Driver Bus 1
  Debug Mode:  true
  Log Level:   info
===================================================

📋 Available MIDI Ports:
───────────────────────────────────────────────────
Input Ports (3):
  1. Keystation 49
  2. IAC Driver Bus 1
  3. Network Session

Output Ports (2):
  1. IAC Driver Bus 1
  2. Network Session
───────────────────────────────────────────────────

🎹 Connecting to MIDI input: Keystation 49
✅ Connected to: Keystation 49

🎵 Listening for MIDI input... (Press Ctrl+C to stop)
```

### Example MIDI Events:

When you play notes or move controls on your MIDI device:

```
[14:23:15.234] 🎹 Note ON  | Ch: 0 | Note:  60 (C4 ) | Vel: 100
[14:23:15.456] 🎹 Note OFF | Ch: 0 | Note:  60 (C4 )
[14:23:16.123] 🎹 Note ON  | Ch: 0 | Note:  64 (E4 ) | Vel:  85
[14:23:16.234] 🎹 Note ON  | Ch: 0 | Note:  67 (G4 ) | Vel:  92
[14:23:16.789] 🎹 Note OFF | Ch: 0 | Note:  64 (E4 )
[14:23:16.890] 🎹 Note OFF | Ch: 0 | Note:  67 (G4 )
[14:23:20.456] 🎛️  CC       | Ch: 0 | CC:   1 | Val:  64
[14:23:21.123] 🎚️  Bend     | Ch: 0 | Value:   2048
[14:23:22.456] 🎼 Program  | Ch: 0 | Program:   5
```

## Code Structure

### Main Components:

**`main.go`:**
- `main()` - Entry point, setup, device connection
- `handleMIDIMessage()` - Processes all incoming MIDI events
- `listMIDIPorts()` - Displays available MIDI ports
- `runWithoutMIDI()` - Fallback mode when no devices available
- `midiNoteToName()` - Converts MIDI note numbers to names (60 → "C4")

### MIDI Message Types Handled:

| Event Type | Symbol | Information Displayed |
|------------|--------|----------------------|
| Note On | 🎹 | Channel, Note Number, Note Name, Velocity |
| Note Off | 🎹 | Channel, Note Number, Note Name |
| Control Change | 🎛️  | Channel, CC Number, Value |
| Pitch Bend | 🎚️  | Channel, Bend Value (-8192 to +8191) |
| Program Change | 🎼 | Channel, Program Number |
| Aftertouch | 👆 | Channel, Pressure |
| Poly Aftertouch | 👆 | Channel, Note, Pressure |
| SysEx | 🔧 | Hexadecimal data dump |

## Configuration

### Environment Variables (.env):

```env
# Your MIDI input device name (must match exactly)
MIDI_INPUT_DEVICE=SE25

# Your MIDI output device name
MIDI_OUTPUT_DEVICE=IAC Driver Bus 1

# Enable debug mode
DEBUG_MODE=true

# Log level
LOG_LEVEL=info
```

### Finding Your MIDI Device Names:

Run the application once to see available devices:

```bash
go run .
```

Look for the "📋 Available MIDI Ports" section in the output.

## Debugging MIDI Events

### Set Breakpoints:

Good locations for debugging in `main.go`:

- **Line 100**: Inside `handleMIDIMessage()` - catches all MIDI events
- **Line 103**: Note On handler - debug note events
- **Line 113**: Control Change handler - debug CC messages
- **Line 40**: Device connection - debug connection issues

### VS Code Debugging:

1. Start Docker: `make debug`
2. Connect debugger: Run > Start Debugging > "Connect to Docker Delve"
3. Set breakpoint in `handleMIDIMessage()`
4. Play notes on MIDI device (if native) or trigger MIDI events
5. Debugger will pause, showing all MIDI message details

## Docker vs Native Development

### Docker (Current Setup):
- ✅ Hot reload works
- ✅ Debugging works
- ❌ No MIDI device access (macOS limitation)
- ✅ Good for developing non-MIDI code
- ✅ Good for CI/CD and testing

### Native (Recommended for MIDI):
- ✅ Full MIDI device access
- ✅ Test with real hardware
- ✅ All environment variables work
- ✅ Hot reload via Air (install: `go install github.com/air-verse/air@latest`)
- ✅ Debugging via Delve

**Recommendation:** Use Docker for general development, run natively when testing MIDI functionality.

## Note Filtering

The `MessageHandler` supports configurable note filtering to accept only specific notes.

### Filter by MIDI Note Numbers

```go
handler := midi.NewMessageHandler()

// Only accept notes C4, D4, E4, F4, G4 (MIDI numbers 60, 62, 64, 65, 67)
handler.SetAllowedNotes([]uint8{60, 62, 64, 65, 67})
```

### Filter by Note Names

```go
handler := midi.NewMessageHandler()

// Only accept notes by name
handler.SetAllowedNoteNames([]string{"C4", "D4", "E4", "F4", "G4"})
```

### Clear Filter (Allow All Notes)

```go
// Reset to allow all notes
handler.SetAllowedNotes(nil)
// or
handler.SetAllowedNoteNames(nil)
```

### Filtered Events

The note filter applies to:
- Note On events
- Note Off events
- Polyphonic Aftertouch events

All other MIDI events (CC, pitch bend, program change, etc.) pass through unfiltered.

## Next Steps

Now that basic MIDI input is working, you can:

1. **Process MIDI events** - Add logic to transform incoming notes
2. **Implement algorithms** - Create composable MIDI processing algorithms
3. **Add MIDI output** - Send processed events to output devices
4. **Build clock system** - Implement timing/sequencing (Phase 2)
5. **Create filters** - Gate certain notes or channels
6. **Add pattern detection** - Recognize and respond to note sequences

## Troubleshooting

### "Could not find MIDI input device"

**In Docker:** Expected - MIDI devices not accessible
**Native:** Check device name exactly matches (case-sensitive)

**Solution:**
```bash
# Run app to see available devices
go run .

# Update .env with exact name from "Available MIDI Ports" list
```

### "build constraints exclude all Go files in rtmidi"

**Cause:** CGO not enabled

**Solution:** Already fixed in docker-compose.yml (`CGO_ENABLED=1`)

### No MIDI events appearing

1. Verify device is connected and powered on
2. Check device is not in use by another application
3. Verify correct device name in `.env`
4. Try running with `DEBUG_MODE=true`

## Implementation Details

### Dependencies:

- `gitlab.com/gomidi/midi/v2` - MIDI library
- `rtmididrv` - Cross-platform MIDI driver (requires CGO)

### Docker Build Changes:

- Added `build-base`, `alsa-lib-dev`, `jack-dev` to Alpine
- Enabled `CGO_ENABLED=1` for MIDI support
- Image size: ~615MB (development)

### Files Modified:

- `main.go` - Complete MIDI implementation
- `Dockerfile` - Added MIDI build dependencies
- `docker-compose.yml` - Enabled CGO
- `go.mod`, `go.sum` - Added MIDI dependencies

---

**Status:** MIDI input handling is fully implemented and tested. Ready for native development with real MIDI devices!
