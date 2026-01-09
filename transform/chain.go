package transform

// Chain executes a sequence of transformers in order
// Output of transformer N becomes input to transformer N+1
type Chain struct {
	transformers []Transformer
	randomness   int // global randomness applied to all transformers
}

// NewChain creates a new transformation chain
func NewChain(randomness int, transformers ...Transformer) *Chain {
	return &Chain{
		transformers: transformers,
		randomness:   randomness,
	}
}

// Process runs the input through all transformers in sequence
func (c *Chain) Process(input MIDINote) []MIDINote {
	if len(c.transformers) == 0 {
		return []MIDINote{input}
	}

	// Start with the input note
	notes := []MIDINote{input}

	// Apply each transformer in sequence
	for _, t := range c.transformers {
		var nextNotes []MIDINote

		// Each transformer processes all notes from the previous stage
		for _, note := range notes {
			result := t.Transform(note, c.randomness)
			nextNotes = append(nextNotes, result...)
		}

		notes = nextNotes
	}

	return notes
}

// ProcessBatch processes multiple input notes through the chain
func (c *Chain) ProcessBatch(inputs []MIDINote) []MIDINote {
	var result []MIDINote
	for _, input := range inputs {
		result = append(result, c.Process(input)...)
	}
	return result
}

// Add appends a transformer to the chain
func (c *Chain) Add(t Transformer) *Chain {
	c.transformers = append(c.transformers, t)
	return c
}

// SetRandomness updates the randomness level (0-10)
func (c *Chain) SetRandomness(r int) *Chain {
	if r < 0 {
		r = 0
	}
	if r > 10 {
		r = 10
	}
	c.randomness = r
	return c
}

// Len returns the number of transformers in the chain
func (c *Chain) Len() int {
	return len(c.transformers)
}

// Names returns the names of all transformers in order
func (c *Chain) Names() []string {
	names := make([]string, len(c.transformers))
	for i, t := range c.transformers {
		names[i] = t.Name()
	}
	return names
}
