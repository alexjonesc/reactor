package midi

import "testing"

func TestSetAllowedNotes(t *testing.T) {
	h := NewMessageHandler()

	// Initially all notes should be allowed
	if !h.isNoteAllowed(60) {
		t.Error("Expected note 60 to be allowed when no filter set")
	}

	// Set allowed notes
	h.SetAllowedNotes([]uint8{60, 64, 67})

	// Check allowed notes pass
	if !h.isNoteAllowed(60) {
		t.Error("Expected note 60 to be allowed")
	}
	if !h.isNoteAllowed(64) {
		t.Error("Expected note 64 to be allowed")
	}
	if !h.isNoteAllowed(67) {
		t.Error("Expected note 67 to be allowed")
	}

	// Check non-allowed notes are blocked
	if h.isNoteAllowed(61) {
		t.Error("Expected note 61 to be blocked")
	}
	if h.isNoteAllowed(0) {
		t.Error("Expected note 0 to be blocked")
	}
	if h.isNoteAllowed(127) {
		t.Error("Expected note 127 to be blocked")
	}
}

func TestSetAllowedNotes_Empty(t *testing.T) {
	h := NewMessageHandler()

	// Set a filter first
	h.SetAllowedNotes([]uint8{60})
	if h.isNoteAllowed(61) {
		t.Error("Expected note 61 to be blocked")
	}

	// Clear with empty slice
	h.SetAllowedNotes([]uint8{})

	// All notes should be allowed again
	if !h.isNoteAllowed(60) {
		t.Error("Expected note 60 to be allowed after clearing filter")
	}
	if !h.isNoteAllowed(61) {
		t.Error("Expected note 61 to be allowed after clearing filter")
	}
}

func TestSetAllowedNotes_Nil(t *testing.T) {
	h := NewMessageHandler()

	// Set a filter first
	h.SetAllowedNotes([]uint8{60})

	// Clear with nil
	h.SetAllowedNotes(nil)

	// All notes should be allowed
	if !h.isNoteAllowed(60) {
		t.Error("Expected note 60 to be allowed after nil reset")
	}
	if !h.isNoteAllowed(127) {
		t.Error("Expected note 127 to be allowed after nil reset")
	}
}

func TestSetAllowedNoteNames(t *testing.T) {
	h := NewMessageHandler()

	// Set allowed note names
	h.SetAllowedNoteNames([]string{"C4", "E4", "G4"})

	// C4 = 60, E4 = 64, G4 = 67
	if !h.isNoteAllowed(60) {
		t.Error("Expected C4 (60) to be allowed")
	}
	if !h.isNoteAllowed(64) {
		t.Error("Expected E4 (64) to be allowed")
	}
	if !h.isNoteAllowed(67) {
		t.Error("Expected G4 (67) to be allowed")
	}

	// D4 = 62 should be blocked
	if h.isNoteAllowed(62) {
		t.Error("Expected D4 (62) to be blocked")
	}
}

func TestSetAllowedNoteNames_Sharps(t *testing.T) {
	h := NewMessageHandler()

	// Set allowed note names with sharps
	h.SetAllowedNoteNames([]string{"C#4", "F#3"})

	// C#4 = 61, F#3 = 54
	if !h.isNoteAllowed(61) {
		t.Error("Expected C#4 (61) to be allowed")
	}
	if !h.isNoteAllowed(54) {
		t.Error("Expected F#3 (54) to be allowed")
	}

	// C4 = 60 should be blocked
	if h.isNoteAllowed(60) {
		t.Error("Expected C4 (60) to be blocked")
	}
}

func TestSetAllowedNoteNames_Empty(t *testing.T) {
	h := NewMessageHandler()

	// Set a filter first
	h.SetAllowedNoteNames([]string{"C4"})

	// Clear with empty slice
	h.SetAllowedNoteNames([]string{})

	// All notes should be allowed
	if !h.isNoteAllowed(60) {
		t.Error("Expected note 60 to be allowed after clearing filter")
	}
	if !h.isNoteAllowed(61) {
		t.Error("Expected note 61 to be allowed after clearing filter")
	}
}

func TestSetAllowedNoteNames_InvalidNames(t *testing.T) {
	h := NewMessageHandler()

	// Set with some invalid names - they should be ignored
	h.SetAllowedNoteNames([]string{"C4", "invalid", "X9", "E4"})

	// Valid notes should work
	if !h.isNoteAllowed(60) {
		t.Error("Expected C4 (60) to be allowed")
	}
	if !h.isNoteAllowed(64) {
		t.Error("Expected E4 (64) to be allowed")
	}

	// Other notes should be blocked
	if h.isNoteAllowed(62) {
		t.Error("Expected D4 (62) to be blocked")
	}
}

func TestSetAllowedNoteNames_NoteZero(t *testing.T) {
	h := NewMessageHandler()

	// C-1 is MIDI note 0
	h.SetAllowedNoteNames([]string{"C-1", "C4"})

	if !h.isNoteAllowed(0) {
		t.Error("Expected C-1 (0) to be allowed")
	}
	if !h.isNoteAllowed(60) {
		t.Error("Expected C4 (60) to be allowed")
	}
}

func TestIsNoteAllowed_DefaultAllowsAll(t *testing.T) {
	h := NewMessageHandler()

	// Test full MIDI range
	for note := uint8(0); note < 128; note++ {
		if !h.isNoteAllowed(note) {
			t.Errorf("Expected note %d to be allowed by default", note)
		}
	}
}
