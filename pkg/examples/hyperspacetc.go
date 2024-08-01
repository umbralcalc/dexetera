package examples

import (
	"math"

	"github.com/umbralcalc/stochadex/pkg/simulator"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Planned approach:
// - Essentially a fake car-following model for the underlying dynamics
// of spacecraft.
// - Use the histogram node iteraton when constructing the node
// controller logic.

// SpacecraftLaneCountStateValueIndices is a mapping which helps with
// describing the meaning of the values for each spacecraft lane count
// state index.
var SpacecraftLaneCountStateValueIndices = map[string]int{
	"Upstream Entry Detection":            0,
	"Downstream Exit Detection":           1,
	"Downstream Queue Size":               2,
	"Latest Upstream Entry Time In Queue": 3,
}

// InverseGaussianSampler - TODO: probably needs to actually be a CDF
type InverseGaussianSampler struct {
	Mu              float64
	Lambda          float64
	unitUniformDist *distuv.Uniform
	unitNormalDist  *distuv.Normal
}

// Rand generates a random sample from an inverse-Gaussian distribution.
func (i *InverseGaussianSampler) Rand() float64 {
	y := math.Pow(i.unitNormalDist.Rand(), 2.0)
	x := i.Mu + (i.Mu * i.Mu * y / (2.0 * i.Lambda)) +
		((i.Mu / (2.0 * i.Lambda)) * math.Sqrt(
			(4.0*i.Mu*i.Lambda*y)+(i.Mu*i.Mu*math.Pow(y, 2.0))))
	if i.unitUniformDist.Rand() <= i.Mu/(i.Mu+x) {
		return x
	} else {
		return i.Mu * i.Mu / x
	}
}

// NewInverseGaussianSampler creates a new InverseGaussianSampler.
func NewInverseGaussianSampler(
	mu float64,
	lambda float64,
	seed uint64,
) *InverseGaussianSampler {
	return &InverseGaussianSampler{
		Mu:     mu,
		Lambda: lambda,
		unitUniformDist: &distuv.Uniform{
			Min: 0.0,
			Max: 1.0,
			Src: rand.NewSource(seed),
		},
		unitNormalDist: &distuv.Normal{
			Mu:    0.0,
			Sigma: 1.0,
			Src:   rand.NewSource(seed),
		},
	}
}

// SpacecraftLaneCountIteration
type SpacecraftLaneCountIteration struct {
	invGaussSampler *InverseGaussianSampler
}

func (s *SpacecraftLaneCountIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
	s.invGaussSampler = NewInverseGaussianSampler(
		0.0,
		1.0,
		settings.Seeds[partitionIndex],
	)
}

func (s *SpacecraftLaneCountIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	// TODO: get lane length parameter
	// TODO: get spacecraft lane speed parameter
	// TODO: use this parameters to set the Mu and Lambda
	return make([]float64, 0)
}
