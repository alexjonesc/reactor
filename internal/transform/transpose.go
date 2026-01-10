package transform

import (
	"math/rand"
)

// Transpose shifts notes by a fixed number of semitones
type Transpose struct {
	Semitones int // positive = up, negative = down
}

func NewTranspose(semitones int) *Transpose {
	return &Transpose{Semitones: semitones}
}

func (t *Transpose) Name() string {
	return "Transpose"
}

func (t *Transpose) Transform(input MIDINote, randomness int) []MIDINote {
	shift := t.Semitones

	// Apply randomness: at max randomness, shift varies by +/- 12 semitones
	if randomness > 0 {
		maxVariation := (randomness * 12) / 10 // scale to max +/- 12
		variation := rand.Intn(maxVariation*2+1) - maxVariation
		shift += variation
	}

	newKey := clampNote(int(input.Key) + shift)

	return []MIDINote{
		{
			Key:      newKey,
			Velocity: input.Velocity,
			Channel:  input.Channel,
		},
	}
}
