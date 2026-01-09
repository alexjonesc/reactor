package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	fmt.Println("🎵 Reactor MIDI Processing System Starting...")
	fmt.Println("=" + repeatString("=", 50))

	// Ensure MIDI driver is closed on exit
	defer midi.CloseDriver()

	// Get configuration from environment variables
	inputDeviceName := getEnv("MIDI_INPUT_DEVICE", "SE25")
	outputDeviceName := getEnv("MIDI_OUTPUT_DEVICE", "IAC Driver Bus 1")
	debugMode := getEnv("DEBUG_MODE", "true")
	logLevel := getEnv("LOG_LEVEL", "info")

	// Display configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  MIDI Input:  %s\n", inputDeviceName)
	fmt.Printf("  MIDI Output: %s\n", outputDeviceName)
	fmt.Printf("  Debug Mode:  %s\n", debugMode)
	fmt.Printf("  Log Level:   %s\n", logLevel)
	fmt.Println("=" + repeatString("=", 50))

	// List available MIDI ports
	listMIDIPorts()

	// Find and open MIDI input device
	fmt.Printf("\n🎹 Connecting to MIDI input: %s\n", inputDeviceName)
	inputPort, err := midi.FindInPort(inputDeviceName)
	if err != nil {
		fmt.Printf("❌ Error: Could not find MIDI input device '%s'\n", inputDeviceName)
		fmt.Printf("   %v\n", err)
		fmt.Println("\nℹ️  Running without MIDI input. Check available ports above.")
		runWithoutMIDI()
		return
	}

	fmt.Printf("✅ Connected to: %s\n", inputPort)

	// Start listening for MIDI messages
	stop, err := midi.ListenTo(inputPort, handleMIDIMessage, midi.UseSysEx())
	if err != nil {
		fmt.Printf("❌ Error listening to MIDI input: %v\n", err)
		return
	}
	defer stop()

	fmt.Println("\n🎵 Listening for MIDI input... (Press Ctrl+C to stop)\n")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Status heartbeat (every 30 seconds)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	heartbeatCount := 0

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n👋 Shutting down gracefully...")
			return
		case <-ticker.C:
			heartbeatCount++
			fmt.Printf("[%s] Status: Listening for MIDI input (heartbeat #%d)\n",
				time.Now().Format("15:04:05"), heartbeatCount)
		}
	}
}

// handleMIDIMessage processes incoming MIDI messages
func handleMIDIMessage(msg midi.Message, timestampms int32) {
	var (
		bt           []byte
		ch, key, vel uint8
		bendLSB      int16
		bendMSB      uint16
	)

	timestamp := time.Now().Format("15:04:05.000")

	switch {
	case msg.GetSysEx(&bt):
		// System Exclusive message
		fmt.Printf("[%s] 🔧 SysEx: % X\n", timestamp, bt)

	case msg.GetNoteStart(&ch, &key, &vel):
		// Note On event
		noteName := midiNoteToName(key)
		fmt.Printf("[%s] 🎹 Note ON  | Ch:%2d | Note: %3d (%s) | Vel: %3d\n",
			timestamp, ch, key, noteName, vel)

	case msg.GetNoteEnd(&ch, &key):
		// Note Off event
		noteName := midiNoteToName(key)
		fmt.Printf("[%s] 🎹 Note OFF | Ch:%2d | Note: %3d (%s)\n",
			timestamp, ch, key, noteName)

	case msg.GetControlChange(&ch, &key, &vel):
		// Control Change (CC) message
		fmt.Printf("[%s] 🎛️  CC       | Ch:%2d | CC: %3d | Val: %3d\n",
			timestamp, ch, key, vel)

	case msg.GetPitchBend(&ch, &bendLSB, &bendMSB):
		// Pitch Bend message
		bend := int(bendLSB) + (int(bendMSB) << 7) - 8192 // Convert to -8192 to +8191
		fmt.Printf("[%s] 🎚️  Bend     | Ch:%2d | Value: %6d\n",
			timestamp, ch, bend)

	case msg.GetProgramChange(&ch, &key):
		// Program Change message
		fmt.Printf("[%s] 🎼 Program  | Ch:%2d | Program: %3d\n",
			timestamp, ch, key)

	case msg.GetAfterTouch(&ch, &vel):
		// Channel Aftertouch
		fmt.Printf("[%s] 👆 AfterTch | Ch:%2d | Pressure: %3d\n",
			timestamp, ch, vel)

	case msg.GetPolyAfterTouch(&ch, &key, &vel):
		// Polyphonic Aftertouch
		noteName := midiNoteToName(key)
		fmt.Printf("[%s] 👆 PolyAT   | Ch:%2d | Note: %3d (%s) | Pressure: %3d\n",
			timestamp, ch, key, noteName, vel)

	default:
		// Unknown or unsupported message
		fmt.Printf("[%s] ❓ Unknown: %v\n", timestamp, msg)
	}
}

// listMIDIPorts displays all available MIDI input and output ports
func listMIDIPorts() {
	fmt.Println("\n📋 Available MIDI Ports:")
	fmt.Println("─" + repeatString("─", 50))

	// List input ports
	inPorts := midi.GetInPorts()
	fmt.Printf("Input Ports (%d):\n", len(inPorts))
	if len(inPorts) == 0 {
		fmt.Println("  (none)")
	} else {
		for i, port := range inPorts {
			fmt.Printf("  %d. %s\n", i+1, port.String())
		}
	}

	fmt.Println()

	// List output ports
	outPorts := midi.GetOutPorts()
	fmt.Printf("Output Ports (%d):\n", len(outPorts))
	if len(outPorts) == 0 {
		fmt.Println("  (none)")
	} else {
		for i, port := range outPorts {
			fmt.Printf("  %d. %s\n", i+1, port.String())
		}
	}

	fmt.Println("─" + repeatString("─", 50))
}

// runWithoutMIDI runs the application without MIDI input (fallback mode)
func runWithoutMIDI() {
	fmt.Println("\n⚠️  Running in fallback mode (no MIDI input)")
	fmt.Println("Press Ctrl+C to stop\n")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		fmt.Printf("[%s] Waiting for MIDI device... (%d)\n",
			time.Now().Format("15:04:05"), count)
	}
}

// midiNoteToName converts a MIDI note number to note name (e.g., 60 -> "C4")
func midiNoteToName(note uint8) string {
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	octave := int(note)/12 - 1
	noteName := noteNames[note%12]
	return fmt.Sprintf("%s%d", noteName, octave)
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// repeatString repeats a string n times
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
