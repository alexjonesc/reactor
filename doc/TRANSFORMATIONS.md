# MIDI Transformation System

A composable system for transforming MIDI notes through chainable algorithms.

## Overview

The transformation system processes MIDI notes through a pipeline of small, single-purpose functions. Each transformer takes a note and produces one or more output notes, allowing complex musical transformations to be built from simple building blocks.

## Core Types

### MIDINote

```go
type MIDINote struct {
    Key      uint8 // MIDI note number (0-127)
    Velocity uint8 // Note velocity (0-127)
    Channel  uint8 // MIDI channel (0-15)
}
```

### Transformer Interface

```go
type Transformer interface {
    Transform(input MIDINote, randomness int) []MIDINote
    Name() string
}
```

- **input**: The MIDI note to transform
- **randomness**: 0-10 scale where 0 = deterministic, 10 = fully random
- **returns**: One or more transformed notes

## Available Transformers

### Transpose

Shifts notes by a fixed number of semitones.

```go
t := transform.NewTranspose(5)  // up 5 semitones
t := transform.NewTranspose(-12) // down one octave
```

| Randomness Effect | Description |
|-------------------|-------------|
| 0 | Exact transposition |
| 1-10 | Pitch varies by up to +/- 12 semitones at max |

### OctaveShift

Adds the input note in additional octaves (note-to-sequence).

```go
o := transform.NewOctaveShift(1)      // add octave above
o := transform.NewOctaveShift(-1, 1)  // add octaves below and above
o := transform.NewOctaveShift(-2, -1, 1, 2) // two octaves each direction
```

| Randomness Effect | Description |
|-------------------|-------------|
| 0 | All specified octaves included |
| 1-10 | Octaves randomly omitted, velocity varies |

### Arpeggio

Transforms a single note into chord tones (note-to-sequence).

```go
a := transform.NewMajorArpeggio()  // major triad
a := transform.NewMinorArpeggio()  // minor triad
a := transform.NewArpeggio([]int{0, 7, 12}) // custom: root, fifth, octave
```

#### Built-in Chord Patterns

| Pattern | Intervals | Notes (from C) |
|---------|-----------|----------------|
| `MajorTriad` | 0, 4, 7 | C, E, G |
| `MinorTriad` | 0, 3, 7 | C, Eb, G |
| `Major7` | 0, 4, 7, 11 | C, E, G, B |
| `Minor7` | 0, 3, 7, 10 | C, Eb, G, Bb |
| `Dominant7` | 0, 4, 7, 10 | C, E, G, Bb |
| `Diminished` | 0, 3, 6 | C, Eb, Gb |
| `Augmented` | 0, 4, 8 | C, E, G# |
| `Sus2` | 0, 2, 7 | C, D, G |
| `Sus4` | 0, 5, 7 | C, F, G |
| `PowerChord` | 0, 7 | C, G |
| `Octaves` | 0, 12 | C, C+oct |
| `FifthStack` | 0, 7, 14 | C, G, D+oct |

| Randomness Effect | Description |
|-------------------|-------------|
| 0 | Exact intervals |
| 1-10 | Intervals vary slightly, velocity varies |

## Chain Executor

The `Chain` pipes multiple transformers together. Output of transformer N becomes input to transformer N+1.

### Creating a Chain

```go
// Method 1: Pass transformers to constructor
chain := transform.NewChain(0,  // randomness level
    transform.NewTranspose(2),
    transform.NewMajorArpeggio(),
)

// Method 2: Build incrementally
chain := transform.NewChain(0)
chain.Add(transform.NewTranspose(2))
chain.Add(transform.NewMajorArpeggio())
```

### Processing Notes

```go
// Single note
input := transform.MIDINote{Key: 60, Velocity: 100, Channel: 0}
output := chain.Process(input)

// Multiple notes
inputs := []transform.MIDINote{
    {Key: 60, Velocity: 100, Channel: 0},
    {Key: 64, Velocity: 90, Channel: 0},
}
output := chain.ProcessBatch(inputs)
```

### Chain Methods

| Method | Description |
|--------|-------------|
| `Process(note)` | Transform a single note |
| `ProcessBatch(notes)` | Transform multiple notes |
| `Add(transformer)` | Append a transformer |
| `SetRandomness(0-10)` | Update randomness level |
| `Len()` | Number of transformers |
| `Names()` | List transformer names |

## Examples

### Simple Transposition

```go
chain := transform.NewChain(0, transform.NewTranspose(7))

input := transform.MIDINote{Key: 60, Velocity: 100, Channel: 0}  // C4
output := chain.Process(input)
// Result: [{Key: 67, Velocity: 100, Channel: 0}]  // G4
```

### Transpose Then Arpeggiate

```go
chain := transform.NewChain(0,
    transform.NewTranspose(2),     // C -> D
    transform.NewMajorArpeggio(),  // D -> D, F#, A
)

input := transform.MIDINote{Key: 60, Velocity: 100, Channel: 0}
output := chain.Process(input)
// Result: D4 (62), F#4 (66), A4 (69)
```

### Octave Layer Then Transpose

```go
chain := transform.NewChain(0,
    transform.NewOctaveShift(1),  // C4 + C5
    transform.NewTranspose(5),    // F4 + F5
)

input := transform.MIDINote{Key: 60, Velocity: 100, Channel: 0}
output := chain.Process(input)
// Result: F4 (65), F5 (77)
```

### Humanized Arpeggios

```go
chain := transform.NewChain(5,  // medium randomness
    transform.NewArpeggio(transform.Minor7),
)

// Each call produces slightly different velocities
// and occasional interval variations
```

## Randomness Guide

| Level | Effect |
|-------|--------|
| 0 | Fully deterministic, identical output every time |
| 1-3 | Subtle variations, tight humanization |
| 4-6 | Noticeable variations, loose feel |
| 7-9 | Significant randomness, unpredictable |
| 10 | Maximum chaos, experimental |

## Creating Custom Transformers

Implement the `Transformer` interface:

```go
type MyTransformer struct {
    // configuration fields
}

func (m *MyTransformer) Name() string {
    return "MyTransformer"
}

func (m *MyTransformer) Transform(input MIDINote, randomness int) []MIDINote {
    // transformation logic
    return []MIDINote{...}
}
```

Then use in a chain:

```go
chain := transform.NewChain(0,
    &MyTransformer{},
    transform.NewTranspose(5),
)
```
