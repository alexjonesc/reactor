package midi

import (
	"fmt"
	"time"

	gomidi "gitlab.com/gomidi/midi/v2"
)

// MessageHandler handles incoming MIDI messages
type MessageHandler struct {
	OnNoteOn       func(channel, key, velocity uint8)
	OnNoteOff      func(channel, key uint8)
	OnControlChange func(channel, cc, value uint8)
	OnPitchBend    func(channel uint8, value int)
	OnProgramChange func(channel, program uint8)
	OnAfterTouch   func(channel, pressure uint8)
	OnPolyAfterTouch func(channel, key, pressure uint8)
	OnSysEx        func(data []byte)
	Verbose        bool
}

// NewMessageHandler creates a handler with default verbose logging
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{Verbose: true}
}

// Handle processes incoming MIDI messages
func (h *MessageHandler) Handle(msg gomidi.Message, timestampms int32) {
	var (
		bt           []byte
		ch, key, vel uint8
		bendLSB      int16
		bendMSB      uint16
	)

	timestamp := time.Now().Format("15:04:05.000")

	switch {
	case msg.GetSysEx(&bt):
		if h.Verbose {
			fmt.Printf("[%s] 🔧 SysEx: % X\n", timestamp, bt)
		}
		if h.OnSysEx != nil {
			h.OnSysEx(bt)
		}

	case msg.GetNoteStart(&ch, &key, &vel):
		if h.Verbose {
			noteName := NoteToName(key)
			fmt.Printf("[%s] 🎹 Note ON  | Ch:%2d | Note: %3d (%s) | Vel: %3d\n",
				timestamp, ch, key, noteName, vel)
		}
		if h.OnNoteOn != nil {
			h.OnNoteOn(ch, key, vel)
		}

	case msg.GetNoteEnd(&ch, &key):
		if h.Verbose {
			noteName := NoteToName(key)
			fmt.Printf("[%s] 🎹 Note OFF | Ch:%2d | Note: %3d (%s)\n",
				timestamp, ch, key, noteName)
		}
		if h.OnNoteOff != nil {
			h.OnNoteOff(ch, key)
		}

	case msg.GetControlChange(&ch, &key, &vel):
		if h.Verbose {
			fmt.Printf("[%s] 🎛️  CC       | Ch:%2d | CC: %3d | Val: %3d\n",
				timestamp, ch, key, vel)
		}
		if h.OnControlChange != nil {
			h.OnControlChange(ch, key, vel)
		}

	case msg.GetPitchBend(&ch, &bendLSB, &bendMSB):
		bend := int(bendLSB) + (int(bendMSB) << 7) - 8192
		if h.Verbose {
			fmt.Printf("[%s] 🎚️  Bend     | Ch:%2d | Value: %6d\n",
				timestamp, ch, bend)
		}
		if h.OnPitchBend != nil {
			h.OnPitchBend(ch, bend)
		}

	case msg.GetProgramChange(&ch, &key):
		if h.Verbose {
			fmt.Printf("[%s] 🎼 Program  | Ch:%2d | Program: %3d\n",
				timestamp, ch, key)
		}
		if h.OnProgramChange != nil {
			h.OnProgramChange(ch, key)
		}

	case msg.GetAfterTouch(&ch, &vel):
		if h.Verbose {
			fmt.Printf("[%s] 👆 AfterTch | Ch:%2d | Pressure: %3d\n",
				timestamp, ch, vel)
		}
		if h.OnAfterTouch != nil {
			h.OnAfterTouch(ch, vel)
		}

	case msg.GetPolyAfterTouch(&ch, &key, &vel):
		if h.Verbose {
			noteName := NoteToName(key)
			fmt.Printf("[%s] 👆 PolyAT   | Ch:%2d | Note: %3d (%s) | Pressure: %3d\n",
				timestamp, ch, key, noteName, vel)
		}
		if h.OnPolyAfterTouch != nil {
			h.OnPolyAfterTouch(ch, key, vel)
		}

	default:
		if h.Verbose {
			fmt.Printf("[%s] ❓ Unknown: %v\n", timestamp, msg)
		}
	}
}
