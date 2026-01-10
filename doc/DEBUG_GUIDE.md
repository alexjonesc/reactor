# VS Code Debugger Test Guide

This guide walks you through testing the VS Code debugger connection with Delve.

## Status: ✅ READY

- **Debugger Status**: Listening on `localhost:2345`
- **Connection Test**: TCP connection verified successful
- **Configuration**: `.vscode/launch.json` is correctly configured

## Quick Test Steps

### 1. Start the Development Environment

The debugger is already running if you just executed `make debug`. If not:

```bash
make debug
```

Expected output:
```
Debugger listening on localhost:2345
Connect your IDE to localhost:2345
```

### 2. Open VS Code

Open this project in VS Code:
```bash
code .
```

### 3. Connect the Debugger

**Option A: Using Run Menu**
1. Go to **Run** > **Start Debugging** (or press `F5`)
2. Select **"Connect to Docker Delve"** from the dropdown
3. The debugger should connect and show "Debugging" in the status bar

**Option B: Using Debug Panel**
1. Click the **Run and Debug** icon in the sidebar (bug icon)
2. Click the green play button next to **"Connect to Docker Delve"**
3. Wait for connection (usually 1-2 seconds)

### 4. Set a Breakpoint

1. Open `main.go`
2. Click in the gutter (left margin) on line 35 (inside the `for range ticker.C` loop)
3. A red dot should appear indicating a breakpoint is set

**Recommended breakpoint locations:**
- **Line 35**: Inside the heartbeat loop (hits every 3 seconds)
- **Line 38**: Inside the milestone condition (hits every 15 seconds)
- **Line 14**: At the start of `main()` (only hits once at startup)

### 5. Verify Debugging Works

Once connected with a breakpoint set:

**Expected behavior:**
- The program will pause when the breakpoint is hit
- VS Code will show:
  - Yellow highlight on the current line
  - Variables panel with current values (`counter`, `inputDevice`, etc.)
  - Call stack in the debug panel
  - Debug controls (Continue, Step Over, Step Into, etc.)

**Test the debugger:**
1. Wait for the breakpoint to hit (up to 3 seconds for line 35)
2. Inspect variables in the **Variables** panel
3. Hover over variables in the code to see their values
4. Try stepping through code with **Step Over** (F10)
5. Click **Continue** (F5) to resume execution

### 6. Stop Debugging

When you're done testing:
- Click **Stop** (red square) in the debug toolbar
- Or press `Shift+F5`
- Or run: `make stop`

## Debugger Configuration Details

### Connection Settings

From `.vscode/launch.json`:
```json
{
    "name": "Connect to Docker Delve",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "remotePath": "/app",
    "port": 2345,
    "host": "localhost",
    "showLog": true,
    "trace": "verbose"
}
```

- **remotePath**: `/app` (container working directory)
- **port**: `2345` (Delve default port)
- **showLog**: `true` (helpful for troubleshooting)
- **trace**: `verbose` (detailed connection logs)

### Path Mappings

The debugger automatically maps:
- **Local**: `/Users/alex/Desktop/reactor/`
- **Container**: `/app/`

This allows you to set breakpoints in your local files that will trigger in the container.

## Troubleshooting

### "Failed to continue: Check the debug console for details"

**Cause**: Delve might not be running or port is not accessible

**Solution**:
```bash
make logs  # Check if Delve is running
lsof -i :2345  # Verify port is listening
make stop && make debug  # Restart container
```

### "Could not attach to process"

**Cause**: Connection timeout or wrong port

**Solution**:
1. Verify container is running: `docker ps | grep reactor-dev`
2. Check logs for "API server listening": `make logs`
3. Test connection: `nc -zv localhost 2345`

### Breakpoints show as "Unverified"

**Cause**: Path mapping issue or code not compiled with debug flags

**Solution**:
1. Verify `.air.docker.toml` has `-gcflags='all=-N -l'` ✓ (already configured)
2. Restart container: `make stop && make debug`
3. Wait for rebuild, then reconnect debugger

### Port 2345 already in use

**Cause**: Another debugger or service using port 2345

**Solution**:
```bash
# Find what's using the port
lsof -i :2345

# Kill the process (if safe) or change Delve port
make stop  # Stop our container first
```

## Advanced: Debug Console Commands

When connected, you can use the **Debug Console** in VS Code to execute Delve commands:

```
dlv> help            # Show available commands
dlv> goroutines      # List all goroutines
dlv> stack           # Show call stack
dlv> print counter   # Print variable value
dlv> locals          # Show local variables
```

## Hot Reload + Debugging

The setup supports both simultaneously:

1. **Edit code** → Air detects change → Rebuilds with debug flags
2. **Auto-restart** → Delve launches the new binary
3. **Debugger reconnects** → Breakpoints remain active
4. **Continue debugging** → No manual intervention needed

This is the power of Air + Delve integration! 🚀

## Test Complete? ✓

If you successfully:
- ✅ Connected VS Code debugger to port 2345
- ✅ Set a breakpoint in `main.go`
- ✅ Paused execution at the breakpoint
- ✅ Inspected variables in the debug panel
- ✅ Stepped through code or continued execution

**Then the debugger is fully functional!** You're ready for MIDI application development with full debugging support.

---

**Next Steps**: Start implementing your MIDI processing code in `main.go` with the confidence that you can debug it at any time.
