package midi

import "fmt"

// NoteNames maps semitone offsets to note names
var NoteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// NoteToName converts a MIDI note number to note name (e.g., 60 -> "C4")
func NoteToName(note uint8) string {
	octave := int(note)/12 - 1
	noteName := NoteNames[note%12]
	return fmt.Sprintf("%s%d", noteName, octave)
}

// NameToNote converts a note name to MIDI note number (e.g., "C4" -> 60)
// Returns 0 if the name is invalid
func NameToNote(name string) uint8 {
	if len(name) < 2 {
		return 0
	}

	var noteName string
	var octaveStr string

	if len(name) >= 2 && name[1] == '#' {
		noteName = name[:2]
		octaveStr = name[2:]
	} else {
		noteName = name[:1]
		octaveStr = name[1:]
	}

	noteIndex := -1
	for i, n := range NoteNames {
		if n == noteName {
			noteIndex = i
			break
		}
	}

	if noteIndex == -1 {
		return 0
	}

	var octave int
	fmt.Sscanf(octaveStr, "%d", &octave)

	return uint8((octave+1)*12 + noteIndex)
}
