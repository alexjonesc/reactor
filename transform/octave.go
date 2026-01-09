package transform

import (
	"math/rand"
)

// OctaveShift adds notes in additional octaves (note-to-sequence)
type OctaveShift struct {
	Octaves []int // octaves to add (e.g., [-1, 1] adds notes one octave below and above)
}

func NewOctaveShift(octaves ...int) *OctaveShift {
	return &OctaveShift{Octaves: octaves}
}

func (o *OctaveShift) Name() string {
	return "OctaveShift"
}

func (o *OctaveShift) Transform(input MIDINote, randomness int) []MIDINote {
	result := []MIDINote{input} // always include original

	for _, oct := range o.Octaves {
		// Randomness affects whether this octave is included
		if randomness > 0 {
			threshold := 10 - randomness // higher randomness = lower threshold
			if rand.Intn(10) < threshold {
				continue // skip this octave randomly
			}
		}

		shift := oct * 12
		newKey := int(input.Key) + shift

		// Only add if within valid MIDI range
		if newKey >= 0 && newKey <= 127 {
			// Randomness also affects velocity
			vel := int(input.Velocity)
			if randomness > 0 {
				variation := rand.Intn(randomness*3) - (randomness * 3 / 2)
				vel = int(clampVelocity(vel + variation))
			}

			result = append(result, MIDINote{
				Key:      uint8(newKey),
				Velocity: uint8(vel),
				Channel:  input.Channel,
			})
		}
	}

	return result
}
