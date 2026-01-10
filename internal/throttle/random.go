package throttle

import (
	"math/rand"

	"reactor/internal/transform"
)

// RandomThrottler accepts notes at random intervals influenced by N
// Higher N = longer average intervals between accepted notes
type RandomThrottler struct {
	n            int                 // throttle intensity (0-8)
	takeFirst    bool                // true = emit first note, false = emit last note
	count        int                 // notes since last emit
	targetCount  int                 // random target for next emit
	buffer       *transform.MIDINote // holds last note when takeFirst=false
}

// NewRandomThrottler creates a random interval throttler
// n: throttle intensity 0-8 (0=no throttle, 8=max throttle)
// takeFirst: true=emit when interval reached, false=buffer until interval
func NewRandomThrottler(n int, takeFirst bool) *RandomThrottler {
	r := &RandomThrottler{
		n:         clampN(n),
		takeFirst: takeFirst,
		count:     0,
		buffer:    nil,
	}
	r.targetCount = r.nextInterval()
	return r
}

func (r *RandomThrottler) Mode() ThrottleMode {
	return ModeRandom
}

func (r *RandomThrottler) Reset() {
	r.count = 0
	r.buffer = nil
	r.targetCount = r.nextInterval()
}

// nextInterval calculates a random interval based on N
// Returns value between 1 and maxInterval
// N=0 -> always 1, N=8 -> between 1 and 12
func (r *RandomThrottler) nextInterval() int {
	if r.n == 0 {
		return 1 // no throttling
	}

	// Base interval grows with N
	// N=1: 1-3, N=4: 1-7, N=8: 1-12
	minInterval := 1
	maxInterval := r.n + (r.n / 2) + 2

	return minInterval + rand.Intn(maxInterval-minInterval+1)
}

func (r *RandomThrottler) Process(note transform.MIDINote) *transform.MIDINote {
	// N=0 means no throttling
	if r.n == 0 {
		return &note
	}

	r.count++

	if r.takeFirst {
		// Emit when we reach target interval
		if r.count >= r.targetCount {
			r.count = 0
			r.targetCount = r.nextInterval()
			return &note
		}
		return nil
	}

	// TakeFirst=false: buffer notes, emit when target reached
	r.buffer = &note

	if r.count >= r.targetCount {
		result := r.buffer
		r.count = 0
		r.buffer = nil
		r.targetCount = r.nextInterval()
		return result
	}

	return nil
}

// SetN updates the throttle intensity
func (r *RandomThrottler) SetN(n int) {
	r.n = clampN(n)
	r.Reset()
}

// N returns the current throttle intensity
func (r *RandomThrottler) N() int {
	return r.n
}

// TakeFirst returns whether first note is captured
func (r *RandomThrottler) TakeFirst() bool {
	return r.takeFirst
}

// SetTakeFirst updates the capture mode
func (r *RandomThrottler) SetTakeFirst(takeFirst bool) {
	r.takeFirst = takeFirst
	r.Reset()
}
