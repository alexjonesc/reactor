package throttle

import "reactor/transform"

// ThrottleMode defines the throttling strategy
type ThrottleMode int

const (
	ModeCount  ThrottleMode = iota // Accept 1 note per N inputs
	ModeRandom                     // Accept notes at random intervals
)

// Throttler limits how often MIDI notes pass through
type Throttler interface {
	// Process evaluates a note and returns it if accepted, nil if throttled
	Process(note transform.MIDINote) *transform.MIDINote
	// Reset clears internal state (counters, buffers)
	Reset()
	// Mode returns the throttle mode
	Mode() ThrottleMode
}

// Config holds throttler configuration
type Config struct {
	Mode      ThrottleMode // Count or Random
	N         int          // 0-8 throttle intensity
	TakeFirst bool         // true = first note, false = last note
}

// clampN ensures N stays within valid range (0-8)
func clampN(n int) int {
	if n < 0 {
		return 0
	}
	if n > 8 {
		return 8
	}
	return n
}
