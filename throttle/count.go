package throttle

import "reactor/transform"

// CountThrottler accepts 1 note for every N+1 input notes
// N=0: all notes pass, N=1: every 2nd note, N=8: every 9th note
type CountThrottler struct {
	n         int  // throttle intensity (0-8)
	takeFirst bool // true = emit first note, false = emit last note
	count     int  // current position in window
	buffer    *transform.MIDINote // holds last note when takeFirst=false
}

// NewCountThrottler creates a count-based throttler
// n: throttle intensity 0-8 (0=no throttle, 8=max throttle)
// takeFirst: true=emit first note in window, false=emit last note
func NewCountThrottler(n int, takeFirst bool) *CountThrottler {
	return &CountThrottler{
		n:         clampN(n),
		takeFirst: takeFirst,
		count:     0,
		buffer:    nil,
	}
}

func (c *CountThrottler) Mode() ThrottleMode {
	return ModeCount
}

func (c *CountThrottler) Reset() {
	c.count = 0
	c.buffer = nil
}

// windowSize returns how many notes per accepted note
// N=0 -> window=1 (all pass), N=1 -> window=2, N=8 -> window=9
func (c *CountThrottler) windowSize() int {
	return c.n + 1
}

func (c *CountThrottler) Process(note transform.MIDINote) *transform.MIDINote {
	// N=0 means no throttling
	if c.n == 0 {
		return &note
	}

	window := c.windowSize()
	c.count++

	if c.takeFirst {
		// Emit on first note of each window
		if c.count == 1 {
			return &note
		}
		if c.count >= window {
			c.count = 0 // reset for next window
		}
		return nil
	}

	// TakeFirst=false: buffer notes, emit when window completes
	c.buffer = &note

	if c.count >= window {
		result := c.buffer
		c.count = 0
		c.buffer = nil
		return result
	}

	return nil
}

// SetN updates the throttle intensity
func (c *CountThrottler) SetN(n int) {
	c.n = clampN(n)
	c.Reset()
}

// N returns the current throttle intensity
func (c *CountThrottler) N() int {
	return c.n
}

// TakeFirst returns whether first note is captured
func (c *CountThrottler) TakeFirst() bool {
	return c.takeFirst
}

// SetTakeFirst updates the capture mode
func (c *CountThrottler) SetTakeFirst(takeFirst bool) {
	c.takeFirst = takeFirst
	c.Reset()
}
