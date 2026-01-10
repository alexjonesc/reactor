# VS Code Debugging Tips for Reactor

## Debug Console Commands

### Common Issue: "debuggee is running"

**Problem:** You can't run debug commands while the program is executing.

**Solution:** Pause execution first!

### How to Pause Execution

**Option 1: Set a Breakpoint (Recommended)**

1. Open `main.go`
2. Click in the left margin (gutter) to set a breakpoint
3. Trigger the breakpoint:
   - For MIDI events: Play notes on your MIDI device
   - For heartbeat: Wait 30 seconds
   - For startup: Restart Air (`make stop-native && make debug-native`)

**Option 2: Manual Pause**

1. Connect debugger: "Connect to Native Delve (Air)"
2. Click **Pause** button in debug toolbar (⏸️)
3. Or press `Cmd + F5` (Mac) / `F6` (Windows/Linux)

### Available Commands When Paused

Once paused at a breakpoint:

#### Inspect Variables
```
dlv> print msg              # Print MIDI message
dlv> print ch               # Print channel
dlv> print key              # Print note number
dlv> print vel              # Print velocity
dlv> print timestamp        # Print timestamp
dlv> locals                 # Show all local variables
dlv> args                   # Show function arguments
```

#### Navigate Stack
```
dlv> stack                  # Show call stack
dlv> up                     # Move up stack frame
dlv> down                   # Move down stack frame
dlv> frame 2                # Jump to specific frame
```

#### List Code
```
dlv> list                   # Show code around current line
dlv> list main.main         # Show specific function
dlv> list main.go:100       # Show specific line
```

#### Control Execution
```
dlv> continue               # Resume execution
dlv> next                   # Step over (next line)
dlv> step                   # Step into function
dlv> stepout                # Step out of function
```

#### Breakpoints
```
dlv> breakpoints            # List all breakpoints
dlv> break main.go:100      # Set breakpoint at line
dlv> clear 1                # Remove breakpoint by ID
dlv> clearall               # Remove all breakpoints
```

#### Goroutines
```
dlv> goroutines             # List all goroutines
dlv> goroutine 2            # Switch to goroutine 2
dlv> threads                # Show all threads
```

#### Advanced
```
dlv> whatis msg             # Show type information
dlv> funcs handleMIDI       # Find functions matching pattern
dlv> types midi.Message     # Show type definition
dlv> vars                   # Show package variables
```

## Good Breakpoint Locations

### MIDI Event Processing (Most Useful)
```go
// Line 100 - handleMIDIMessage() entry
func handleMIDIMessage(msg midi.Message, timestampms int32) {
    // Set breakpoint here - triggers on EVERY MIDI event
    var (
        bt           []byte
```

### Note On Events
```go
// Line 100 - Note On handler
case msg.GetNoteStart(&ch, &key, &vel):
    // Set breakpoint here - triggers when notes are played
    noteName := midiNoteToName(key)
```

### Control Change
```go
// Line 112 - CC handler
case msg.GetControlChange(&ch, &key, &vel):
    // Set breakpoint here - triggers on knob/fader movements
    fmt.Printf("[%s] 🎛️  CC       | Ch:%2d | CC: %3d | Val: %3d\n",
```

### Device Connection
```go
// Line 40 - Device connection
inputPort, err := midi.FindInPort(inputDeviceName)
if err != nil {
    // Set breakpoint here - triggers during startup
    fmt.Printf("❌ Error: Could not find MIDI input device '%s'\n", inputDeviceName)
```

### Heartbeat Loop
```go
// Line 77 - Heartbeat (triggers every 30 seconds)
case <-ticker.C:
    heartbeatCount++
    // Set breakpoint here for periodic inspection
    fmt.Printf("[%s] Status: Listening for MIDI input (heartbeat #%d)\n",
```

## Debugging Workflow Examples

### Example 1: Debug Note Events

1. **Set breakpoint at line 100** (Note On handler)
2. **Start debugger**: `make debug-native`
3. **Connect VS Code**: Select "Connect to Native Delve (Air)"
4. **Play a note** on your MIDI keyboard
5. **Debugger pauses** - now inspect:
   ```
   dlv> print ch       # Channel number
   dlv> print key      # Note number (60 = C4)
   dlv> print vel      # Velocity (0-127)
   dlv> print noteName # Note name string
   ```
6. **Step through**: Press F10 to see each line execute
7. **Continue**: Press F5 to resume

### Example 2: Debug Control Changes

1. **Set breakpoint at line 112** (CC handler)
2. **Move a knob/fader** on your MIDI controller
3. **Debugger pauses**:
   ```
   dlv> locals         # See all variables
   dlv> print key      # CC number
   dlv> print vel      # CC value (0-127)
   ```

### Example 3: Watch Variable Values

1. **Set breakpoint at line 100**
2. When paused, add to **Watch** panel:
   - `ch` - See channel in real-time
   - `key` - See note numbers
   - `vel` - See velocity
   - `noteName` - See note names

### Example 4: Conditional Breakpoint

Right-click in gutter → **Add Conditional Breakpoint**

**Break only on middle C:**
```
key == 60
```

**Break only on high velocity:**
```
vel > 100
```

**Break only on channel 1:**
```
ch == 1
```

## VS Code Debug Shortcuts

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Start/Continue | F5 | F5 |
| Pause | Cmd+F5 | F6 |
| Step Over | F10 | F10 |
| Step Into | F11 | F11 |
| Step Out | Shift+F11 | Shift+F11 |
| Restart | Cmd+Shift+F5 | Ctrl+Shift+F5 |
| Stop | Shift+F5 | Shift+F5 |
| Toggle Breakpoint | F9 | F9 |

## Debug Console Tips

### Multi-line Commands

You can't directly run multi-line commands, but you can:

1. Use semicolons:
   ```
   dlv> print ch; print key; print vel
   ```

2. Or run separate commands:
   ```
   dlv> print ch
   dlv> print key
   dlv> print vel
   ```

### Command History

- **Up Arrow** - Previous command
- **Down Arrow** - Next command

### Clear Console

VS Code: Click the trash icon (🗑️) in debug console toolbar

## Hot Reload While Debugging

When Air detects changes:

1. **Air rebuilds** the application
2. **Delve restarts** with the new binary
3. **Breakpoints remain active** (if lines haven't changed)
4. **Debugger reconnects** automatically

This means you can:
- Edit code while debugging
- Let Air rebuild
- Continue debugging immediately!

## Troubleshooting

### "debuggee is running" Error

**Cause:** Program is executing, not paused

**Solutions:**
1. Set a breakpoint and trigger it
2. Click Pause button in debug toolbar
3. Press `Cmd + F5` to pause

### Breakpoint Not Hitting

**Possible reasons:**

1. **Code not executing**
   - Check if your code path is reached
   - Add logging before breakpoint to verify

2. **Wrong file/line**
   - Ensure breakpoint is in the right location
   - Check line numbers match compiled code

3. **Debugger not attached**
   - Verify "Connected" in debug toolbar
   - Check port 2346 is listening: `lsof -i :2346`

4. **Breakpoint in wrong binary**
   - Restart: `make stop-native && make debug-native`
   - Ensure Air compiled with debug flags

### Variables Not Showing

**Cause:** Compiler optimization

**Solution:** Already configured! `.air.native.toml` uses:
```
-gcflags='all=-N -l'
```
This disables optimizations for debugging.

### "could not find symbol value" Error

**Cause:** Variable optimized away or out of scope

**Solutions:**
1. Check variable is in current scope
2. Move up/down stack frames: `up` / `down`
3. Use `locals` to see available variables

## Advanced: Debugging Without VS Code

If you prefer command-line debugging:

### Start Delve Directly

```bash
# Stop Air
make stop-native

# Build with debug flags
go build -gcflags='all=-N -l' -o tmp/main .

# Start Delve
dlv exec --listen=:2346 --api-version=2 ./tmp/main

# In another terminal, connect
dlv connect :2346
```

### Delve Command-Line Interface

```bash
(dlv) break main.handleMIDIMessage
(dlv) continue
(dlv) print msg
(dlv) locals
(dlv) stack
```

## Performance Impact

**Debugging has minimal performance impact:**

- Delve adds ~5-10ms latency
- Acceptable for MIDI (usually <10ms jitter tolerance)
- If you need absolute lowest latency, run without debugger:
  ```bash
  go run .
  ```

## Summary

**To run debug commands:**

1. ✅ **Pause first** - Set breakpoint or click Pause
2. ✅ **Wait for pause** - Program must be stopped
3. ✅ **Run commands** - Now you can inspect variables
4. ✅ **Continue** - Resume execution when done

**Quick Test:**

```bash
# Start debugger
make debug-native

# In VS Code:
# 1. Set breakpoint at main.go:100
# 2. Play a MIDI note
# 3. Debugger pauses
# 4. Run: print key
# 5. Continue (F5)
```

---

**Remember:** The debugger is your friend! Set breakpoints liberally and inspect everything. 🐛🔍
