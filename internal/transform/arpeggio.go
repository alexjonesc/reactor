package transform

import (
	"math/rand"
)

// Chord interval patterns (semitones from root)
var (
	MajorTriad  = []int{0, 4, 7}       // major chord
	MinorTriad  = []int{0, 3, 7}       // minor chord
	Major7      = []int{0, 4, 7, 11}   // major 7th
	Minor7      = []int{0, 3, 7, 10}   // minor 7th
	Dominant7   = []int{0, 4, 7, 10}   // dominant 7th
	Diminished  = []int{0, 3, 6}       // diminished
	Augmented   = []int{0, 4, 8}       // augmented
	Sus2        = []int{0, 2, 7}       // suspended 2nd
	Sus4        = []int{0, 5, 7}       // suspended 4th
	PowerChord  = []int{0, 7}          // power chord (5th)
	Octaves     = []int{0, 12}         // octave
	FifthStack  = []int{0, 7, 14}      // stacked fifths
)

// Arpeggio transforms a single note into chord tones (note-to-sequence)
type Arpeggio struct {
	Intervals []int // semitone intervals from root (e.g., [0, 4, 7] for major)
}

func NewArpeggio(intervals []int) *Arpeggio {
	return &Arpeggio{Intervals: intervals}
}

func NewMajorArpeggio() *Arpeggio {
	return NewArpeggio(MajorTriad)
}

func NewMinorArpeggio() *Arpeggio {
	return NewArpeggio(MinorTriad)
}

func (a *Arpeggio) Name() string {
	return "Arpeggio"
}

func (a *Arpeggio) Transform(input MIDINote, randomness int) []MIDINote {
	result := make([]MIDINote, 0, len(a.Intervals))

	for _, interval := range a.Intervals {
		// Randomness affects interval slightly (microtonal-ish humanization)
		actualInterval := interval
		if randomness > 0 && interval != 0 {
			// At high randomness, intervals can vary by +/- 2 semitones
			maxShift := (randomness * 2) / 10
			if maxShift > 0 {
				actualInterval += rand.Intn(maxShift*2+1) - maxShift
			}
		}

		newKey := int(input.Key) + actualInterval
		if newKey < 0 || newKey > 127 {
			continue
		}

		// Randomness affects velocity
		vel := int(input.Velocity)
		if randomness > 0 {
			variation := rand.Intn(randomness*2+1) - randomness
			vel = int(clampVelocity(vel + variation))
		}

		result = append(result, MIDINote{
			Key:      uint8(newKey),
			Velocity: uint8(vel),
			Channel:  input.Channel,
		})
	}

	return result
}
