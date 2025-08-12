package main

import (
	"fmt"
	"reactor/clock"
	"reactor/sequencer"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {

	// Midi setup ----------------------------------- //
	defer midi.CloseDriver()

	in, err := midi.FindInPort("SE25")
	if err != nil {
		fmt.Println("can't find SE25")
		return
	}

	// Get output ports
	outPorts := midi.GetOutPorts()
	fmt.Println("---MIDI output ports---")
	for i, port := range outPorts {
		fmt.Printf("%d \"%s\"\n", i, port.String())
	}

	out, err := midi.FindOutPort("IAC Driver Bus 1")
	if err != nil {
		fmt.Printf("can't find IAC Driver Bus 1")
		return
	}
	out.Open()
	defer out.Close()

	fmt.Println(out)

	// Sequencer setup ----------------------------- //
	var sq1 = sequencer.New(&out)

	// Clock setup --------------------------------- //
	var clock = clock.New()
	clock.Subscribe(&sq1)

	// Midi Listen ------------------------------- //
	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, note, vel uint8
		//var chh = uint8(15)

		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		case msg.GetNoteStart(&ch, &note, &vel):
			//fmt.Printf("starting note %s on channel %v with velocity %v\n", note, ch, vel)
			midiInput := sequencer.MidiInut{Ch: ch, Note: note, Vel: vel}

			clock.Tick(midiInput)
			sq1.Input(midiInput)
			//sq1.DoNothing(midiInput)
		case msg.GetNoteEnd(&ch, &note):
			//sq1.Stop()
			//fmt.Printf("ending note %s on channel %v\n", midi.Note(note), ch)
			//out.Send(midi.NoteOff(chh, note, vel))
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	time.Sleep(time.Second * 1000)

	stop()
}
