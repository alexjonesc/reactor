package throttle

import (
	"testing"

	"reactor/internal/transform"
)

func makeNote(key uint8) transform.MIDINote {
	return transform.MIDINote{Key: key, Velocity: 100, Channel: 0}
}

// --- CountThrottler Tests ---

func TestCountThrottler_NoThrottle(t *testing.T) {
	th := NewCountThrottler(0, true)

	// All notes should pass through
	for i := 0; i < 10; i++ {
		note := makeNote(uint8(60 + i))
		result := th.Process(note)
		if result == nil {
			t.Errorf("note %d: expected to pass, got nil", i)
		}
		if result != nil && result.Key != note.Key {
			t.Errorf("note %d: expected key %d, got %d", i, note.Key, result.Key)
		}
	}
}

func TestCountThrottler_TakeFirst_N1(t *testing.T) {
	th := NewCountThrottler(1, true) // every 2nd note, take first

	// Input:  0 1 2 3 4 5 6 7
	// Output: 0   2   4   6
	expected := []bool{true, false, true, false, true, false, true, false}

	for i, exp := range expected {
		note := makeNote(uint8(60 + i))
		result := th.Process(note)

		if exp && result == nil {
			t.Errorf("note %d: expected to pass, got nil", i)
		}
		if !exp && result != nil {
			t.Errorf("note %d: expected throttled, got %v", i, result)
		}
	}
}

func TestCountThrottler_TakeFirst_N4(t *testing.T) {
	th := NewCountThrottler(4, true) // every 5th note, take first

	// Input:  0 1 2 3 4 5 6 7 8 9
	// Output: 0         5
	expected := []bool{true, false, false, false, false, true, false, false, false, false}

	for i, exp := range expected {
		note := makeNote(uint8(60 + i))
		result := th.Process(note)

		if exp && result == nil {
			t.Errorf("note %d: expected to pass, got nil", i)
		}
		if !exp && result != nil {
			t.Errorf("note %d: expected throttled, got %v", i, result)
		}
	}
}

func TestCountThrottler_TakeLast_N1(t *testing.T) {
	th := NewCountThrottler(1, false) // every 2nd note, take last

	// Input:  0 1 2 3 4 5 6 7
	// Output:   1   3   5   7
	notes := []transform.MIDINote{
		makeNote(60), makeNote(61), makeNote(62), makeNote(63),
		makeNote(64), makeNote(65), makeNote(66), makeNote(67),
	}
	expectedPass := []bool{false, true, false, true, false, true, false, true}
	expectedKeys := []uint8{0, 61, 0, 63, 0, 65, 0, 67}

	for i, note := range notes {
		result := th.Process(note)

		if expectedPass[i] {
			if result == nil {
				t.Errorf("note %d: expected to pass, got nil", i)
			} else if result.Key != expectedKeys[i] {
				t.Errorf("note %d: expected key %d, got %d", i, expectedKeys[i], result.Key)
			}
		} else {
			if result != nil {
				t.Errorf("note %d: expected throttled, got %v", i, result)
			}
		}
	}
}

func TestCountThrottler_TakeLast_N4(t *testing.T) {
	th := NewCountThrottler(4, false) // every 5th note, take last

	// Send 10 notes, should emit on 5th and 10th
	var emitted []uint8
	for i := 0; i < 10; i++ {
		note := makeNote(uint8(60 + i))
		result := th.Process(note)
		if result != nil {
			emitted = append(emitted, result.Key)
		}
	}

	if len(emitted) != 2 {
		t.Fatalf("expected 2 emitted notes, got %d", len(emitted))
	}
	if emitted[0] != 64 { // 5th note (index 4)
		t.Errorf("expected first emit key 64, got %d", emitted[0])
	}
	if emitted[1] != 69 { // 10th note (index 9)
		t.Errorf("expected second emit key 69, got %d", emitted[1])
	}
}

func TestCountThrottler_Reset(t *testing.T) {
	th := NewCountThrottler(2, true) // every 3rd note

	// Process 2 notes
	th.Process(makeNote(60))
	th.Process(makeNote(61))

	// Reset
	th.Reset()

	// Next note should pass (first in new window)
	result := th.Process(makeNote(62))
	if result == nil {
		t.Error("expected note to pass after reset")
	}
}

func TestCountThrottler_Mode(t *testing.T) {
	th := NewCountThrottler(4, true)
	if th.Mode() != ModeCount {
		t.Errorf("expected ModeCount, got %v", th.Mode())
	}
}

func TestCountThrottler_SetN(t *testing.T) {
	th := NewCountThrottler(2, true)

	th.SetN(5)
	if th.N() != 5 {
		t.Errorf("expected N=5, got %d", th.N())
	}

	// Should clamp to bounds
	th.SetN(-1)
	if th.N() != 0 {
		t.Errorf("expected N=0 after -1, got %d", th.N())
	}

	th.SetN(10)
	if th.N() != 8 {
		t.Errorf("expected N=8 after 10, got %d", th.N())
	}
}

func TestCountThrottler_MaxThrottle(t *testing.T) {
	th := NewCountThrottler(8, true) // max throttle: every 9th note

	passCount := 0
	for i := 0; i < 27; i++ { // 3 windows
		result := th.Process(makeNote(60))
		if result != nil {
			passCount++
		}
	}

	if passCount != 3 {
		t.Errorf("expected 3 notes to pass (27/9), got %d", passCount)
	}
}

// --- RandomThrottler Tests ---

func TestRandomThrottler_NoThrottle(t *testing.T) {
	th := NewRandomThrottler(0, true)

	// All notes should pass through
	for i := 0; i < 10; i++ {
		note := makeNote(uint8(60 + i))
		result := th.Process(note)
		if result == nil {
			t.Errorf("note %d: expected to pass, got nil", i)
		}
	}
}

func TestRandomThrottler_TakeFirst_Throttles(t *testing.T) {
	th := NewRandomThrottler(4, true)

	// With N=4, we should throttle some notes
	// Run enough notes to see some throttling
	passCount := 0
	totalNotes := 100

	for i := 0; i < totalNotes; i++ {
		result := th.Process(makeNote(60))
		if result != nil {
			passCount++
		}
	}

	// Should pass fewer notes than total
	if passCount >= totalNotes {
		t.Errorf("expected some throttling, got %d/%d passed", passCount, totalNotes)
	}
	// But should pass at least some
	if passCount == 0 {
		t.Error("expected some notes to pass, got 0")
	}
}

func TestRandomThrottler_TakeLast_BuffersCorrectly(t *testing.T) {
	th := NewRandomThrottler(2, false)

	// Send notes with different keys
	var emittedKeys []uint8
	for i := 0; i < 50; i++ {
		note := makeNote(uint8(60 + (i % 12)))
		result := th.Process(note)
		if result != nil {
			emittedKeys = append(emittedKeys, result.Key)
		}
	}

	// Should have emitted some notes
	if len(emittedKeys) == 0 {
		t.Error("expected some notes to be emitted")
	}
}

func TestRandomThrottler_HigherN_LessNotes(t *testing.T) {
	// Higher N should result in fewer notes passing through

	countAtN2 := 0
	th2 := NewRandomThrottler(2, true)
	for i := 0; i < 1000; i++ {
		if th2.Process(makeNote(60)) != nil {
			countAtN2++
		}
	}

	countAtN8 := 0
	th8 := NewRandomThrottler(8, true)
	for i := 0; i < 1000; i++ {
		if th8.Process(makeNote(60)) != nil {
			countAtN8++
		}
	}

	if countAtN8 >= countAtN2 {
		t.Errorf("expected N=8 to pass fewer notes than N=2, got N=2:%d, N=8:%d", countAtN2, countAtN8)
	}
}

func TestRandomThrottler_Reset(t *testing.T) {
	th := NewRandomThrottler(4, false)

	// Process some notes
	for i := 0; i < 5; i++ {
		th.Process(makeNote(60))
	}

	// Reset
	th.Reset()

	// Internal state should be cleared
	// Can't easily verify random state, but should not panic
	th.Process(makeNote(61))
}

func TestRandomThrottler_Mode(t *testing.T) {
	th := NewRandomThrottler(4, true)
	if th.Mode() != ModeRandom {
		t.Errorf("expected ModeRandom, got %v", th.Mode())
	}
}

func TestRandomThrottler_SetN(t *testing.T) {
	th := NewRandomThrottler(2, true)

	th.SetN(6)
	if th.N() != 6 {
		t.Errorf("expected N=6, got %d", th.N())
	}

	// Should clamp to bounds
	th.SetN(-5)
	if th.N() != 0 {
		t.Errorf("expected N=0 after -5, got %d", th.N())
	}

	th.SetN(15)
	if th.N() != 8 {
		t.Errorf("expected N=8 after 15, got %d", th.N())
	}
}

// --- Helper Function Tests ---

func TestClampN(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{-5, 0},
		{0, 0},
		{4, 4},
		{8, 8},
		{10, 8},
	}

	for _, tt := range tests {
		result := clampN(tt.input)
		if result != tt.expected {
			t.Errorf("clampN(%d): expected %d, got %d", tt.input, tt.expected, result)
		}
	}
}

// --- Interface Compliance Tests ---

func TestThrottlerInterface(t *testing.T) {
	// Verify both throttlers implement the interface
	var _ Throttler = &CountThrottler{}
	var _ Throttler = &RandomThrottler{}
}

// --- Integration-style Tests ---

func TestCountThrottler_StreamSimulation(t *testing.T) {
	th := NewCountThrottler(3, true) // every 4th note (window=4)

	// Simulate a stream of MIDI notes
	input := []uint8{60, 62, 64, 65, 67, 69, 71, 72, 74, 76, 77, 79}
	var output []uint8

	for _, key := range input {
		result := th.Process(makeNote(key))
		if result != nil {
			output = append(output, result.Key)
		}
	}

	// Window size = 4, so 12 notes = 3 windows
	// First note of each window: indices 0, 4, 8
	expected := []uint8{60, 67, 74}
	if len(output) != len(expected) {
		t.Fatalf("expected %d notes, got %d: %v", len(expected), len(output), output)
	}
	for i, exp := range expected {
		if output[i] != exp {
			t.Errorf("output[%d]: expected %d, got %d", i, exp, output[i])
		}
	}
}

func TestCountThrottler_TakeLast_StreamSimulation(t *testing.T) {
	th := NewCountThrottler(3, false) // every 4th note, take last

	input := []uint8{60, 62, 64, 65, 67, 69, 71, 72, 74, 76, 77, 79}
	var output []uint8

	for _, key := range input {
		result := th.Process(makeNote(key))
		if result != nil {
			output = append(output, result.Key)
		}
	}

	// Should get last note of each 4-note window
	expected := []uint8{65, 72, 79}
	if len(output) != len(expected) {
		t.Fatalf("expected %d notes, got %d: %v", len(expected), len(output), output)
	}
	for i, exp := range expected {
		if output[i] != exp {
			t.Errorf("output[%d]: expected %d, got %d", i, exp, output[i])
		}
	}
}
