# MIDI Throttling System

A configurable system for limiting how often MIDI notes pass through, useful for thinning dense note streams or creating rhythmic gates.

## Overview

The throttling system filters incoming MIDI notes based on configurable rules. It supports both deterministic (count-based) and probabilistic (random) throttling modes, with options to capture either the first or last note in each throttle window.

## Core Types

### ThrottleMode

```go
type ThrottleMode int

const (
    ModeCount  ThrottleMode = iota  // Deterministic: 1 note per N inputs
    ModeRandom                       // Probabilistic: random intervals
)
```

### Throttler Interface

```go
type Throttler interface {
    Process(note MIDINote) *MIDINote  // Returns note if accepted, nil if throttled
    Reset()                            // Clear internal state
    Mode() ThrottleMode               // Return throttle mode
}
```

### Config

```go
type Config struct {
    Mode      ThrottleMode  // Count or Random
    N         int           // 0-8 throttle intensity
    TakeFirst bool          // true = first note, false = last note
}
```

## Throttle Intensity (N)

The `N` parameter controls throttle intensity on a 0-8 scale:

| N | Effect |
|---|--------|
| 0 | No throttling - all notes pass through |
| 1-3 | Light throttling |
| 4-5 | Medium throttling |
| 6-7 | Heavy throttling |
| 8 | Maximum throttling |

## Throttle Modes

### Count-Based Throttling

Accepts exactly 1 note for every N+1 input notes. Deterministic and predictable.

```go
th := throttle.NewCountThrottler(n int, takeFirst bool)
```

**Window Size Formula:** `window = N + 1`

| N | Window | Notes Passed |
|---|--------|--------------|
| 0 | 1 | All (no throttle) |
| 1 | 2 | Every 2nd note |
| 2 | 3 | Every 3rd note |
| 3 | 4 | Every 4th note |
| 4 | 5 | Every 5th note |
| 8 | 9 | Every 9th note |

**Visual Example (N=3, window=4):**

```
TakeFirst=true:
Input:  ● ● ● ● ● ● ● ● ● ● ● ●
Output: ●       ●       ●        (notes 1, 5, 9)

TakeFirst=false:
Input:  ● ● ● ● ● ● ● ● ● ● ● ●
Output:       ●       ●       ●  (notes 4, 8, 12)
```

### Random Interval Throttling

Accepts notes at random intervals influenced by N. Higher N means longer average intervals.

```go
th := throttle.NewRandomThrottler(n int, takeFirst bool)
```

**Interval Range Formula:** `minInterval=1, maxInterval=N + N/2 + 2`

| N | Interval Range | Behavior |
|---|----------------|----------|
| 0 | Always 1 | All notes pass |
| 2 | 1-5 | Light random throttle |
| 4 | 1-8 | Medium random throttle |
| 8 | 1-14 | Heavy random throttle |

**Visual Example (N=3, randomized):**

```
Input:  ● ● ● ● ● ● ● ● ● ● ● ●
Output: ●   ●     ●   ●     ●    (varies each time)
```

## TakeFirst vs TakeLast

Controls which note is captured within each throttle window.

### TakeFirst (true)

Emits the **first** note when a new window begins. Immediate response, no latency.

```
Window:  [1 2 3 4] [5 6 7 8]
Emitted:  ^         ^
```

**Use cases:**
- Real-time performance where immediate response matters
- Triggering samples on first hit
- Gate effects

### TakeLast (false)

Buffers notes and emits the **last** note when the window completes. Adds latency but captures the most recent value.

```
Window:  [1 2 3 4] [5 6 7 8]
Emitted:        ^         ^
```

**Use cases:**
- Capturing final position of rapid movements
- Smoothing rapid pitch bends or mod wheel
- When you want the "settled" value

## Usage Examples

### Basic Count Throttling

```go
// Pass every 4th note, capture first
th := throttle.NewCountThrottler(3, true)

for _, note := range incomingNotes {
    if result := th.Process(note); result != nil {
        fmt.Printf("Note passed: %d\n", result.Key)
    }
}
```

### Random Throttling for Humanization

```go
// Random intervals, medium intensity
th := throttle.NewRandomThrottler(4, true)

for _, note := range incomingNotes {
    if result := th.Process(note); result != nil {
        // Notes pass at irregular intervals
        playNote(*result)
    }
}
```

### Dynamic Throttle Adjustment

```go
th := throttle.NewCountThrottler(2, true)

// Increase throttling based on note density
if notesPerSecond > 100 {
    th.SetN(6)  // Heavy throttle
} else if notesPerSecond > 50 {
    th.SetN(3)  // Medium throttle
} else {
    th.SetN(0)  // No throttle
}
```

### Switching Capture Mode

```go
th := throttle.NewCountThrottler(4, true)

// Switch to capturing last note
th.SetTakeFirst(false)

// Check current mode
if th.TakeFirst() {
    fmt.Println("Capturing first note")
} else {
    fmt.Println("Capturing last note")
}
```

### Reset State

```go
th := throttle.NewCountThrottler(3, true)

// Process some notes...
th.Process(note1)
th.Process(note2)

// Reset clears counters and buffers
th.Reset()

// Next note starts fresh window
```

## API Reference

### CountThrottler

| Method | Description |
|--------|-------------|
| `NewCountThrottler(n, takeFirst)` | Create count-based throttler |
| `Process(note) *MIDINote` | Process note, return if accepted |
| `Reset()` | Clear state |
| `Mode() ThrottleMode` | Returns `ModeCount` |
| `N() int` | Get current N value |
| `SetN(n)` | Set N value (clamped 0-8) |
| `TakeFirst() bool` | Get capture mode |
| `SetTakeFirst(bool)` | Set capture mode |

### RandomThrottler

| Method | Description |
|--------|-------------|
| `NewRandomThrottler(n, takeFirst)` | Create random throttler |
| `Process(note) *MIDINote` | Process note, return if accepted |
| `Reset()` | Clear state, generate new interval |
| `Mode() ThrottleMode` | Returns `ModeRandom` |
| `N() int` | Get current N value |
| `SetN(n)` | Set N value (clamped 0-8) |
| `TakeFirst() bool` | Get capture mode |
| `SetTakeFirst(bool)` | Set capture mode |

## Comparison: Count vs Random

| Aspect | Count | Random |
|--------|-------|--------|
| Predictability | Deterministic | Variable |
| Rhythm | Regular pattern | Irregular |
| Use case | Rhythmic gates | Humanization |
| Latency (TakeFirst) | None | None |
| Output rate | Exact (1/window) | Approximate |

## Performance Considerations

- Both throttlers are stateful and maintain minimal internal state
- `Process()` is O(1) - suitable for real-time MIDI processing
- No allocations during normal operation (returns pointer to input or nil)
- Thread-safety: Not thread-safe by default; wrap with mutex if needed for concurrent access
