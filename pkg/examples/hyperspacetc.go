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

// InverseGaussianSampler
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

// SpacecraftFollowingLaneIteration
type SpacecraftFollowingLaneIteration struct {
}

func (s *SpacecraftFollowingLaneIteration) Configure(
	partitionIndex int,
	settings *simulator.Settings,
) {
}

func (s *SpacecraftFollowingLaneIteration) Iterate(
	params *simulator.OtherParams,
	partitionIndex int,
	stateHistories []*simulator.StateHistory,
	timestepsHistory *simulator.CumulativeTimestepsHistory,
) []float64 {
	return make([]float64, 0)
}
