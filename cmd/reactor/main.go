package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reactor/internal/midi"

	gomidi "gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	fmt.Println("Reactor MIDI Processing System Starting...")
	fmt.Println("=" + repeatString("=", 50))

	// Ensure MIDI driver is closed on exit
	defer gomidi.CloseDriver()

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
	midi.ListPorts()

	// Find and open MIDI input device
	fmt.Printf("\nConnecting to MIDI input: %s\n", inputDeviceName)
	inputPort, err := gomidi.FindInPort(inputDeviceName)
	if err != nil {
		fmt.Printf("Error: Could not find MIDI input device '%s'\n", inputDeviceName)
		fmt.Printf("   %v\n", err)
		fmt.Println("\nRunning without MIDI input. Check available ports above.")
		runWithoutMIDI()
		return
	}

	fmt.Printf("Connected to: %s\n", inputPort)

	// Create message handler
	handler := midi.NewMessageHandler()

	// Only allow B3 and C4 notes
	handler.SetAllowedNoteNames([]string{"B3", "C4"})

	// Start listening for MIDI messages
	stop, err := gomidi.ListenTo(inputPort, handler.Handle, gomidi.UseSysEx())
	if err != nil {
		fmt.Printf("Error listening to MIDI input: %v\n", err)
		return
	}
	defer stop()

	fmt.Println("\nListening for MIDI input... (Press Ctrl+C to stop)")

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
			fmt.Println("\n\nShutting down gracefully...")
			return
		case <-ticker.C:
			heartbeatCount++
			fmt.Printf("[%s] Status: Listening for MIDI input (heartbeat #%d)\n",
				time.Now().Format("15:04:05"), heartbeatCount)
		}
	}
}

// runWithoutMIDI runs the application without MIDI input (fallback mode)
func runWithoutMIDI() {
	fmt.Println("\nRunning in fallback mode (no MIDI input)")
	fmt.Println("Press Ctrl+C to stop")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		fmt.Printf("[%s] Waiting for MIDI device... (%d)\n",
			time.Now().Format("15:04:05"), count)
	}
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
