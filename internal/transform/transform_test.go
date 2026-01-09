package transform

import (
	"testing"
)

// --- Transpose Tests ---

func TestTranspose_Up(t *testing.T) {
	tr := NewTranspose(5) // up 5 semitones
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0} // C4

	result := tr.Transform(input, 0) // no randomness

	if len(result) != 1 {
		t.Fatalf("expected 1 note, got %d", len(result))
	}
	if result[0].Key != 65 { // F4
		t.Errorf("expected key 65, got %d", result[0].Key)
	}
	if result[0].Velocity != 100 {
		t.Errorf("expected velocity 100, got %d", result[0].Velocity)
	}
}

func TestTranspose_Down(t *testing.T) {
	tr := NewTranspose(-12) // down one octave
	input := MIDINote{Key: 72, Velocity: 80, Channel: 1} // C5

	result := tr.Transform(input, 0)

	if result[0].Key != 60 { // C4
		t.Errorf("expected key 60, got %d", result[0].Key)
	}
	if result[0].Channel != 1 {
		t.Errorf("expected channel 1, got %d", result[0].Channel)
	}
}

func TestTranspose_ClampLow(t *testing.T) {
	tr := NewTranspose(-100) // way down
	input := MIDINote{Key: 10, Velocity: 100, Channel: 0}

	result := tr.Transform(input, 0)

	if result[0].Key != 0 { // should clamp to 0
		t.Errorf("expected key 0, got %d", result[0].Key)
	}
}

func TestTranspose_ClampHigh(t *testing.T) {
	tr := NewTranspose(100) // way up
	input := MIDINote{Key: 100, Velocity: 100, Channel: 0}

	result := tr.Transform(input, 0)

	if result[0].Key != 127 { // should clamp to 127
		t.Errorf("expected key 127, got %d", result[0].Key)
	}
}

func TestTranspose_Name(t *testing.T) {
	tr := NewTranspose(5)
	if tr.Name() != "Transpose" {
		t.Errorf("expected name 'Transpose', got '%s'", tr.Name())
	}
}

// --- OctaveShift Tests ---

func TestOctaveShift_SingleOctaveUp(t *testing.T) {
	os := NewOctaveShift(1) // add one octave up
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}

	result := os.Transform(input, 0)

	if len(result) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(result))
	}
	if result[0].Key != 60 { // original
		t.Errorf("expected original key 60, got %d", result[0].Key)
	}
	if result[1].Key != 72 { // octave up
		t.Errorf("expected shifted key 72, got %d", result[1].Key)
	}
}

func TestOctaveShift_MultipleOctaves(t *testing.T) {
	os := NewOctaveShift(-1, 1) // add octave below and above
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}

	result := os.Transform(input, 0)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(result))
	}

	keys := make(map[uint8]bool)
	for _, n := range result {
		keys[n.Key] = true
	}

	if !keys[60] || !keys[48] || !keys[72] {
		t.Errorf("expected keys 48, 60, 72, got %v", keys)
	}
}

func TestOctaveShift_OutOfRange(t *testing.T) {
	os := NewOctaveShift(10) // way up - should be out of range for high notes
	input := MIDINote{Key: 120, Velocity: 100, Channel: 0}

	result := os.Transform(input, 0)

	// Should only have original note since 120 + 120 > 127
	if len(result) != 1 {
		t.Errorf("expected 1 note (original only), got %d", len(result))
	}
}

func TestOctaveShift_Name(t *testing.T) {
	os := NewOctaveShift(1)
	if os.Name() != "OctaveShift" {
		t.Errorf("expected name 'OctaveShift', got '%s'", os.Name())
	}
}

// --- Arpeggio Tests ---

func TestArpeggio_MajorTriad(t *testing.T) {
	arp := NewMajorArpeggio()
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0} // C4

	result := arp.Transform(input, 0)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes for major triad, got %d", len(result))
	}

	expected := []uint8{60, 64, 67} // C, E, G
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
	}
}

func TestArpeggio_MinorTriad(t *testing.T) {
	arp := NewMinorArpeggio()
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0} // C4

	result := arp.Transform(input, 0)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes for minor triad, got %d", len(result))
	}

	expected := []uint8{60, 63, 67} // C, Eb, G
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
	}
}

func TestArpeggio_CustomIntervals(t *testing.T) {
	arp := NewArpeggio([]int{0, 7, 12}) // root, fifth, octave
	input := MIDINote{Key: 48, Velocity: 90, Channel: 2}

	result := arp.Transform(input, 0)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(result))
	}

	expected := []uint8{48, 55, 60}
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
		if result[i].Channel != 2 {
			t.Errorf("note %d: expected channel 2, got %d", i, result[i].Channel)
		}
	}
}

func TestArpeggio_PreservesVelocity(t *testing.T) {
	arp := NewMajorArpeggio()
	input := MIDINote{Key: 60, Velocity: 77, Channel: 0}

	result := arp.Transform(input, 0)

	for i, note := range result {
		if note.Velocity != 77 {
			t.Errorf("note %d: expected velocity 77, got %d", i, note.Velocity)
		}
	}
}

func TestArpeggio_Name(t *testing.T) {
	arp := NewMajorArpeggio()
	if arp.Name() != "Arpeggio" {
		t.Errorf("expected name 'Arpeggio', got '%s'", arp.Name())
	}
}

// --- Chain Tests ---

func TestChain_Empty(t *testing.T) {
	chain := NewChain(0) // no transformers
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}

	result := chain.Process(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 note, got %d", len(result))
	}
	if result[0].Key != 60 {
		t.Errorf("expected key 60, got %d", result[0].Key)
	}
}

func TestChain_SingleTransformer(t *testing.T) {
	chain := NewChain(0, NewTranspose(12)) // octave up
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}

	result := chain.Process(input)

	if len(result) != 1 {
		t.Fatalf("expected 1 note, got %d", len(result))
	}
	if result[0].Key != 72 {
		t.Errorf("expected key 72, got %d", result[0].Key)
	}
}

func TestChain_MultipleTransformers(t *testing.T) {
	// Transpose up 2, then create major arpeggio
	chain := NewChain(0,
		NewTranspose(2),    // D4 from C4
		NewMajorArpeggio(), // D, F#, A
	)
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0} // C4

	result := chain.Process(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(result))
	}

	expected := []uint8{62, 66, 69} // D4, F#4, A4
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
	}
}

func TestChain_ExpansionThenTransform(t *testing.T) {
	// First expand to octaves, then transpose all
	chain := NewChain(0,
		NewOctaveShift(1),  // original + octave up
		NewTranspose(5),    // transpose all up a fourth
	)
	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}

	result := chain.Process(input)

	if len(result) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(result))
	}

	expected := []uint8{65, 77} // F4, F5
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
	}
}

func TestChain_Add(t *testing.T) {
	chain := NewChain(0)
	chain.Add(NewTranspose(5))
	chain.Add(NewTranspose(7))

	if chain.Len() != 2 {
		t.Errorf("expected 2 transformers, got %d", chain.Len())
	}

	input := MIDINote{Key: 60, Velocity: 100, Channel: 0}
	result := chain.Process(input)

	if result[0].Key != 72 { // 60 + 5 + 7
		t.Errorf("expected key 72, got %d", result[0].Key)
	}
}

func TestChain_ProcessBatch(t *testing.T) {
	chain := NewChain(0, NewTranspose(12))
	inputs := []MIDINote{
		{Key: 60, Velocity: 100, Channel: 0},
		{Key: 64, Velocity: 90, Channel: 0},
		{Key: 67, Velocity: 80, Channel: 0},
	}

	result := chain.ProcessBatch(inputs)

	if len(result) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(result))
	}

	expected := []uint8{72, 76, 79}
	for i, exp := range expected {
		if result[i].Key != exp {
			t.Errorf("note %d: expected key %d, got %d", i, exp, result[i].Key)
		}
	}
}

func TestChain_Names(t *testing.T) {
	chain := NewChain(0,
		NewTranspose(5),
		NewOctaveShift(1),
		NewMajorArpeggio(),
	)

	names := chain.Names()

	expected := []string{"Transpose", "OctaveShift", "Arpeggio"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d names, got %d", len(expected), len(names))
	}
	for i, exp := range expected {
		if names[i] != exp {
			t.Errorf("name %d: expected '%s', got '%s'", i, exp, names[i])
		}
	}
}

func TestChain_SetRandomness(t *testing.T) {
	chain := NewChain(5)
	chain.SetRandomness(8)
	// Can't easily test internal state, but we can verify it doesn't panic

	chain.SetRandomness(-1) // should clamp to 0
	chain.SetRandomness(15) // should clamp to 10
}

// --- Helper Function Tests ---

func TestClampNote(t *testing.T) {
	tests := []struct {
		input    int
		expected uint8
	}{
		{-10, 0},
		{0, 0},
		{60, 60},
		{127, 127},
		{200, 127},
	}

	for _, tt := range tests {
		result := clampNote(tt.input)
		if result != tt.expected {
			t.Errorf("clampNote(%d): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}

func TestClampVelocity(t *testing.T) {
	tests := []struct {
		input    int
		expected uint8
	}{
		{-5, 0},
		{0, 0},
		{100, 100},
		{127, 127},
		{150, 127},
	}

	for _, tt := range tests {
		result := clampVelocity(tt.input)
		if result != tt.expected {
			t.Errorf("clampVelocity(%d): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}
