package rating

import "math"

type Rating struct {
	Mu  float64
	Sig float64
}

func (r *Rating) WithNoise(sigNoise float64) Rating {
	return Rating{
		Mu:  r.Mu,
		Sig: math.Hypot(r.Sig, sigNoise),
	}
}

func (r *Rating) TowardsNoise(decay float64, limit *Rating) Rating {
	muDiff := r.Mu - limit.Mu
	siSqDiff := r.Sig*r.Sig - limit.Sig*limit.Sig
	return Rating{
		Mu:  limit.Mu + muDiff*decay,
		Sig: math.Sqrt(limit.Sig*limit.Sig + siSqDiff*decay*decay),
	}
}
