package traffic

import (
	"fmt"
	"log"
	"time"

	"github.com/dop251/goja"
	"golang.org/x/exp/rand"
)

// Source represents a source of events in a given time range.
type Source struct {
	// From gives the starting tick (remainder) for this source.
	//
	// Defaults to zero.
	//
	// This value is sampled at each new System.Width tick.
	From *Dist

	from int64

	// To gives the last tick (remainder, exclusive) for this source.
	//
	// Defaults to System.Width.
	//
	// This value is sampled at each new System.Width tick.
	To *Dist

	to int64

	// D sample times D.Scale gives the number of events per tick.
	D *Dist

	// JS is optional Javascript code that should result in the
	// number of events for a tick.
	//
	// The variable 't' and 'r', which is 't' mod the system
	// Width, are in the Javascript environment when this string
	// of code is evaluated.
	//
	// The environment persists through the run.
	JS string

	vm *goja.Runtime
}

func (s *Source) Reset(width, t, r int64) {
	if s.From != nil {
		s.from = int64(s.From.Rand())
	}
	if s.To == nil {
		s.to = width
	} else {
		s.to = int64(s.To.Rand())
	}
}

func (s *Source) Count(t, r int64) int64 {
	if s.D != nil {
		n := int64(s.D.Scale * s.D.Rand())
		if n < 0 {
			return 0
		}
		return n
	}

	s.vm.Set("t", t)
	s.vm.Set("r", r)
	v, err := s.vm.RunString(s.JS)
	if err != nil {
		panic(err)
	}

	var n int64
	x := v.Export()

	switch vv := x.(type) {
	case int64:
		n = vv
	case float64:
		n = int64(vv)
	default:
		panic(fmt.Errorf("code returned %#v (%T)", vv, vv))
	}

	if n < 0 {
		return 0
	}

	return n
}

// System is a set of event sources.
type System struct {
	// Sources maps arbitrary, opaque names to Sources, each of
	// which can contribute to the aggregate count returned by
	// Counts().
	Sources map[string]*Source

	// Width is the modulus for the ticker.
	//
	// Defaults to 60.
	//
	// This value is effectively the largest tick, after which the
	// clock resets to zero and continues.
	Width int64

	// Log turns on some logging output.
	Log bool
}

// Init validates distributions and initializes RNGs.
func (s *System) Init(r rand.Source) error {
	if r == nil {
		r = rand.NewSource(uint64(time.Now().UnixNano()))
	}

	if s.Width == 0 {
		s.Width = 60
	}

	for name, src := range s.Sources {
		wrap := func(err error) error {
			return fmt.Errorf("in '%s': %v", name, err)
		}
		if src.From != nil {
			if err := src.From.Validate(); err != nil {
				return wrap(err)
			}
			src.From.SetSrc(r)
		}
		if src.To != nil {
			if err := src.To.Validate(); err != nil {
				return wrap(err)
			}
			src.To.SetSrc(r)
		}

		if src.D != nil {
			if err := src.D.Validate(); err != nil {
				return wrap(err)
			}
			if src.D.Scale == 0 {
				src.D.Scale = 1
			}
			src.D.SetSrc(r)
		}

		if src.D == nil && src.JS == "" {
			return wrap(fmt.Errorf("No D or JS"))
		}

		if src.D != nil && src.JS != "" {
			return wrap(fmt.Errorf("Can't have both N and JS"))
		}

		if src.JS != "" {
			src.vm = goja.New()
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
		if r == 0 {
			d.Reset(s.Width, t, r)
		}
		if r < d.from || d.to <= r {
			continue
		}

		n := d.Count(t, r)
		if s.Log {
			log.Printf("traffic %d %s %d", t, name, n)
		}
		counts[name] += n
	}
	return counts
}
