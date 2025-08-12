package clock

import (
	"fmt"
	"reactor/helpers"
	"reactor/sequencer"
)

type Clock struct {
	tick        *int
	subdiv      *int
	beat        *int
	beatSubdiv  *int
	subscribers []*sequencer.Sequencer
}

func New() Clock {
	var (
		tick        = 0
		subdiv      = 4 // per beat
		beat        = 0
		beatSubdiv  = 4
		subscribers = []*sequencer.Sequencer{}
	)

	var clock = Clock{
		tick:        &tick,
		subdiv:      &subdiv,
		beat:        &beat,
		beatSubdiv:  &beatSubdiv,
		subscribers: subscribers,
	}

	return clock
}

func (c *Clock) Tick(input sequencer.MidiInut) {
	*c.tick = helpers.Mod(*c.tick, *c.subdiv) + 1

	if *c.tick == 1 {
		*c.beat = helpers.Mod(*c.beat, *c.beatSubdiv) + 1
	}
	c.Publish()
	fmt.Println("Clock:", *c.beat, *c.tick)
}

// Figure out this logic later!
// func Mod(a int, b int) int {
// 	return a - a/b*b
// }

func (c *Clock) Subscribe(subscriber *sequencer.Sequencer) {
	c.subscribers = append(c.subscribers, subscriber)
}

func (c *Clock) Publish() {
	for _, subscriber := range c.subscribers {
		subscriber.Notify(*c.tick)
	}
}
