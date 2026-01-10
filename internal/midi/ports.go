package midi

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi/v2"
)

// ListPorts displays all available MIDI input and output ports
func ListPorts() {
	fmt.Println("\n📋 Available MIDI Ports:")
	fmt.Println("─" + repeatString("─", 50))

	// List input ports
	inPorts := gomidi.GetInPorts()
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
	outPorts := gomidi.GetOutPorts()
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

// GetInputPorts returns all available MIDI input ports
func GetInputPorts() []string {
	ports := gomidi.GetInPorts()
	names := make([]string, len(ports))
	for i, p := range ports {
		names[i] = p.String()
	}
	return names
}

// GetOutputPorts returns all available MIDI output ports
func GetOutputPorts() []string {
	ports := gomidi.GetOutPorts()
	names := make([]string, len(ports))
	for i, p := range ports {
		names[i] = p.String()
	}
	return names
}

func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
