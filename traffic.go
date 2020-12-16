package traffic

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"
)

// Source represents a source of events in a given time range.
type Source struct {
	// From gives the starting tick (remainder) for this source.
	From Dist

	// To gives the last tick (remainder, exclusive) for this source.
	To Dist

	// Scale is multiplied by N.Rand() to give the number of
	// events for tick.
	Scale float64

	// N returns the number of events per 1/Scale.
	N Dist
}

// System is a set of event sources.
type System struct {
	// Sources maps arbitrary, opaque names to Sources, each of
	// which can contribute to the aggregate count returned by
	// Counts().
	Sources map[string]*Source

	// Width is the modulus for the ticker.
	//
	// This value is effectively the largest tick, after which the
	// clock resets to zero and continues.
	Width int64

	// NoWarn turns off some possible warnings.
	NoWarn bool
}

// Init validates distributions and initializes RNGs.
func (s *System) Init(r rand.Source) error {
	if r == nil {
		r = rand.NewSource(uint64(time.Now().UnixNano()))
	}
	for name, src := range s.Sources {
		wrap := func(err error) error {
			return fmt.Errorf("in '%s': %v", name, err)
		}
		if err := src.From.Validate(); err != nil {
			return wrap(err)
		}
		src.From.SetSrc(r)
		if err := src.To.Validate(); err != nil {
			return wrap(err)
		}
		src.To.SetSrc(r)
		if err := src.N.Validate(); err != nil {
			return wrap(err)
		}
		src.N.SetSrc(r)

		if src.Scale == 0 {
			if !s.NoWarn {
				fmt.Fprintf(os.Stderr, "warning: %s Scale is 0", name)
			}
		}
	}

	return nil
}

// Counts is the primary method, which returns the number of events by
// source.
func (s *System) Counts(t int64) map[string]int64 {
	var (
		r      = t % s.Width
		counts = make(map[string]int64, len(s.Sources))
	)
	for name, d := range s.Sources {
		var (
			from = int64(d.From.Rand())
			to   = int64(d.To.Rand())
		)
		if r < from || to <= r {
			continue
		}
		n := int64(d.Scale * d.N.Rand())
		if 0 < n {
			counts[name] += n
		}
	}
	return counts
}
