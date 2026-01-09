package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("🎵 Reactor MIDI Processing System Starting... [HOT RELOAD TEST]")
	fmt.Println("=" + repeatString("=", 50))

	// Get configuration from environment variables
	inputDevice := getEnv("MIDI_INPUT_DEVICE", "SE25")
	outputDevice := getEnv("MIDI_OUTPUT_DEVICE", "IAC Driver Bus 1")
	debugMode := getEnv("DEBUG_MODE", "false")
	logLevel := getEnv("LOG_LEVEL", "info")

	// Display configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  MIDI Input:  %s\n", inputDevice)
	fmt.Printf("  MIDI Output: %s\n", outputDevice)
	fmt.Printf("  Debug Mode:  %s\n", debugMode)
	fmt.Printf("  Log Level:   %s\n", logLevel)
	fmt.Println("=" + repeatString("=", 50))

	// Simple heartbeat to show the app is running
	fmt.Println("\nReactor is running... (Hot reload is active)")
	fmt.Println("Press Ctrl+C to stop\n")

	counter := 0
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		counter++
		// 🔍 DEBUGGER TEST: Set a breakpoint on the next line to test VS Code debugging
		fmt.Printf("[%s] Heartbeat #%d - Ready for MIDI input from %s\n",
			time.Now().Format("15:04:05"), counter, inputDevice)

		// Add a simple condition to demonstrate debugging
		if counter%5 == 0 {
			fmt.Printf("  → Milestone: %d heartbeats completed\n", counter)
		}
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
