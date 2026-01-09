package transform

// MIDINote represents a single MIDI note with all relevant properties
type MIDINote struct {
	Key      uint8 // MIDI note number (0-127)
	Velocity uint8 // Note velocity (0-127)
	Channel  uint8 // MIDI channel (0-15)
}

// Transformer transforms MIDI notes with optional randomization
// randomness: 0 = deterministic, 10 = fully random
type Transformer interface {
	Transform(input MIDINote, randomness int) []MIDINote
	Name() string
}

// clampNote ensures a note value stays within valid MIDI range (0-127)
func clampNote(n int) uint8 {
	if n < 0 {
		return 0
	}
	if n > 127 {
		return 127
	}
	return uint8(n)
}

// clampVelocity ensures velocity stays within valid range (0-127)
func clampVelocity(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 127 {
		return 127
	}
	return uint8(v)
}
