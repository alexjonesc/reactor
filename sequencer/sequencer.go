package sequencer

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

type Gate struct {
	count     *uint8
	threshold *uint8
}

type ActiveNote struct {
	note   uint8
	length uint8
	tick   *uint8
}

type Sequencer struct {
	out         *drivers.Out
	keeper      map[uint8]Gate
	activeNotes map[uint8]ActiveNote
	pattern     []uint8
	pointer     *int
}

type MidiInut struct {
	Ch   uint8
	Note uint8
	Vel  uint8
}

var (
	outCh       uint8
	currentTick int
)

var pat1 = []uint8{37, 38, 39, 40}

func init() {
	outCh = uint8(15)
}

// TODO: Figure out how to send a stepped sequence of outgoing midi messages.
func New(out *drivers.Out) Sequencer {
	var (
		keeper    = make(map[uint8]Gate)
		count     = uint8(0)
		threshold = uint8(5)
		pointer   = int(-1)
	)

	keeper[uint8(37)] = Gate{
		count:     &count,
		threshold: &threshold,
	}
	//keeper[uint8(39)] = Gate{count: &count, threshold: &threshold, pattern: pat1}

	var sq = Sequencer{
		out:         out,
		keeper:      keeper,
		activeNotes: make(map[uint8]ActiveNote),
		pattern:     pat1,
		pointer:     &pointer,
	}

	return sq
}

func (sq *Sequencer) inputGate(input MidiInut) bool {
	note := input.Note
	gate, exists := sq.keeper[note]

	// check if note is allowed
	if !exists {
		return false
	}

	// increment counter and check threshold
	*gate.count++
	if *gate.count >= *gate.threshold-1 {
		*gate.count = 0
		return true
	}

	return false
}

func (sq *Sequencer) DoNothing(input MidiInut) {
	//..
}

func (sq *Sequencer) Input(input MidiInut) {

	// input handler
	play := sq.inputGate(input)
	if !play {
		return
	}

	// play
	prevNote, nextNote := sq.prevAndNext()
	sq.stop(prevNote)
	sq.play(input, nextNote)
}

func (sq *Sequencer) prevAndNext() (uint8, uint8) {
	prevPointer := max(*sq.pointer, 0) // because it can be -1 to start
	*sq.pointer = incrementPtr(*sq.pointer, sq.pattern)

	prevNote := sq.pattern[prevPointer]
	nextNote := sq.pattern[*sq.pointer]

	return prevNote, nextNote
}

func incrementPtr(p int, a []uint8) int {
	p++

	if p > len(a)-1 {
		p = 0
	}

	return p
}

func (sq *Sequencer) play(input MidiInut, note uint8) {
	var (
		out       drivers.Out = *sq.out
		vel                   = uint8(127)
		startTick             = uint8(0)
	)

	out.Send(midi.NoteOff(outCh, note))
	err := out.Send(midi.NoteOn(outCh, note, vel))
	sq.activeNotes[note] = ActiveNote{note: note, length: uint8(4), tick: &startTick}

	if err != nil {
		fmt.Println("Error sending NoteOn:", err)
	}

	fmt.Println("-> NoteOn:", note, sq.activeNotes)
}

func (sq *Sequencer) stop(note uint8) {
	var out drivers.Out = *sq.out

	_, exists := sq.activeNotes[note]
	if exists {
		out.Send(midi.NoteOff(outCh, note))
		delete(sq.activeNotes, note)
		fmt.Println("-> NoteOff:", note, sq.activeNotes)
	}
}

func (sq *Sequencer) Notify(tick int) {
	currentTick = tick

	// check active note lengths
	for k, a := range sq.activeNotes {
		*a.tick++
		if *a.tick >= a.length {
			sq.stop(a.note)
		}
		fmt.Println("---->", k, *a.tick)
	}

	// fmt.Println(sq.activeNotes)

}
