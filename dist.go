package traffic

import (
	"fmt"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Dist should represent a single underlying distribution in a manner
// that makes deserialization type-safe.
type Dist struct {
	// Scale, which defaults to 1, is multiplied by Dist.Rand() to
	// give the number of events for a tick.
	Scale float64

	Const        *Const               `json:",omitempty"`
	AlphaStable  *distuv.AlphaStable  `json:",omitempty"`
	Bernoulli    *distuv.Bernoulli    `json:",omitempty"`
	Beta         *distuv.Beta         `json:",omitempty"`
	Binomial     *distuv.Binomial     `json:",omitempty"`
	Categorical  *distuv.Categorical  `json:",omitempty"`
	ChiSquared   *distuv.ChiSquared   `json:",omitempty"`
	Exponential  *distuv.Exponential  `json:",omitempty"`
	F            *distuv.F            `json:",omitempty"`
	Gamma        *distuv.Gamma        `json:",omitempty"`
	GumbelRight  *distuv.GumbelRight  `json:",omitempty"`
	InverseGamma *distuv.InverseGamma `json:",omitempty"`
	Laplace      *distuv.Laplace      `json:",omitempty"`
	LogNormal    *distuv.LogNormal    `json:",omitempty"`
	Normal       *distuv.Normal       `json:",omitempty"`
	Pareto       *distuv.Pareto       `json:",omitempty"`
	Poisson      *distuv.Poisson      `json:",omitempty"`
	StudentsT    *distuv.StudentsT    `json:",omitempty"`
	Triangle     *distuv.Triangle     `json:",omitempty"`
	Uniform      *distuv.Uniform      `json:",omitempty"`
	Weibull      *distuv.Weibull      `json:",omitempty"`
}

// Validate ensures that exactly one underlying distribution is
// present.
func (d *Dist) Validate() error {
	n := 0
	if d.Const != nil {
		n++
	}
	if d.AlphaStable != nil {
		n++
	}
	if d.Bernoulli != nil {
		n++
	}
	if d.Beta != nil {
		n++
	}
	if d.Binomial != nil {
		n++
	}
	if d.Categorical != nil {
		n++
	}
	if d.ChiSquared != nil {
		n++
	}
	if d.Exponential != nil {
		n++
	}
	if d.F != nil {
		n++
	}
	if d.Gamma != nil {
		n++
	}
	if d.GumbelRight != nil {
		n++
	}
	if d.InverseGamma != nil {
		n++
	}
	if d.Laplace != nil {
		n++
	}
	if d.LogNormal != nil {
		n++
	}
	if d.Normal != nil {
		n++
	}
	if d.Pareto != nil {
		n++
	}
	if d.Poisson != nil {
		n++
	}
	if d.StudentsT != nil {
		n++
	}
	if d.Triangle != nil {
		n++
	}
	if d.Uniform != nil {
		n++
	}
	if d.Weibull != nil {
		n++
	}

	if n != 1 {
		return fmt.Errorf("saw %d underlying distributions instead of exactly one", n)
	}

	return nil
}

// SetSrc sets the source for the RNG.
func (d *Dist) SetSrc(r rand.Source) {
	if d.AlphaStable != nil {
		d.AlphaStable.Src = r
	}

	if d.AlphaStable != nil {
		d.AlphaStable.Src = r
	}

	if d.Bernoulli != nil {
		d.Bernoulli.Src = r
	}

	if d.Beta != nil {
		d.Beta.Src = r
	}

	if d.Binomial != nil {
		d.Binomial.Src = r
	}

	if d.ChiSquared != nil {
		d.ChiSquared.Src = r
	}

	if d.Exponential != nil {
		d.Exponential.Src = r
	}

	if d.F != nil {
		d.F.Src = r
	}

	if d.Gamma != nil {
		d.Gamma.Src = r
	}

	if d.GumbelRight != nil {
		d.GumbelRight.Src = r
	}

	if d.InverseGamma != nil {
		d.InverseGamma.Src = r
	}

	if d.Laplace != nil {
		d.Laplace.Src = r
	}

	if d.LogNormal != nil {
		d.LogNormal.Src = r
	}

	if d.Normal != nil {
		d.Normal.Src = r
	}

	if d.Pareto != nil {
		d.Pareto.Src = r
	}

	if d.Poisson != nil {
		d.Poisson.Src = r
	}

	if d.StudentsT != nil {
		d.StudentsT.Src = r
	}

	if d.Uniform != nil {
		d.Uniform.Src = r
	}

	if d.Weibull != nil {
		d.Weibull.Src = r
	}
}

// Rand samples from the (underlying) distribution.
func (d *Dist) Rand() float64 {
	if d.Const != nil {
		return d.Const.Rand()
	}
	if d.AlphaStable != nil {
		return d.AlphaStable.Rand()
	}

	if d.Bernoulli != nil {
		return d.Bernoulli.Rand()
	}

	if d.Beta != nil {
		return d.Beta.Rand()
	}

	if d.Binomial != nil {
		return d.Binomial.Rand()
	}

	if d.Categorical != nil {
		return d.Categorical.Rand()
	}

	if d.ChiSquared != nil {
		return d.ChiSquared.Rand()
	}

	if d.Exponential != nil {
		return d.Exponential.Rand()
	}

	if d.F != nil {
		return d.F.Rand()
	}

	if d.Gamma != nil {
		return d.Gamma.Rand()
	}

	if d.GumbelRight != nil {
		return d.GumbelRight.Rand()
	}

	if d.InverseGamma != nil {
		return d.InverseGamma.Rand()
	}

	if d.Laplace != nil {
		return d.Laplace.Rand()
	}

	if d.LogNormal != nil {
		return d.LogNormal.Rand()
	}

	if d.Normal != nil {
		return d.Normal.Rand()
	}

	if d.Pareto != nil {
		return d.Pareto.Rand()
	}

	if d.Poisson != nil {
		return d.Poisson.Rand()
	}

	if d.StudentsT != nil {
		return d.StudentsT.Rand()
	}

	if d.Triangle != nil {
		return d.Triangle.Rand()
	}

	if d.Uniform != nil {
		return d.Uniform.Rand()
	}

	if d.Weibull != nil {
		return d.Weibull.Rand()
	}

	panic(fmt.Errorf("didn't find underlying distribution"))
}

// Distribution is a univariate distribution.
type Distribution interface {
	Rand() float64
}

// Const is a constant Distribution.
type Const float64

func (d Const) Rand() float64 {
	return float64(d)
}
