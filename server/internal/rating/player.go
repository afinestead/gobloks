package rating

import "math"

type PlayerEvent struct {
	ContestIdx uint
	Place      uint
	RatingMu   int32
	RatingSig  int32
	PerfScore  int32
}

type Player struct {
	NormalFactor    Rating
	LogisticFactors []TanhTerm
	EventHistory    []*PlayerEvent
	ApproxPosterior Rating
	UpdateTime      uint64
	DeltaTime       uint64
}

func (p *Player) withRating(mu, sig float64, updateTime uint64) Player {
	return Player{
		NormalFactor:    Rating{Mu: mu, Sig: sig},
		LogisticFactors: []TanhTerm{},
		EventHistory:    []*PlayerEvent{},
		ApproxPosterior: Rating{Mu: mu, Sig: sig},
		UpdateTime:      updateTime,
		DeltaTime:       0,
	}
}

func (p *Player) timesPlayed() uint {
	return uint(len(p.EventHistory))
}

func (p *Player) timesPlayedExcl() uint {
	return p.timesPlayed() - 1
}

func (p *Player) isNewcomer() bool {
	return p.timesPlayedExcl() == 0
}

func (p *Player) updateRating(rating Rating, performanceScore float64) {
	lastEvent := p.EventHistory[len(p.EventHistory)-1]

	p.ApproxPosterior = rating
	lastEvent.RatingMu = int32(math.Round(rating.Mu))
	lastEvent.RatingSig = int32(math.Round(rating.Sig))
	lastEvent.PerfScore = int32(math.Round(performanceScore))
}

func (p *Player) addNoiseBest(sigNoise, transferSpeed float64) {
	newPosterior := p.ApproxPosterior.WithNoise(sigNoise)
	decay := math.Pow(p.ApproxPosterior.Sig/newPosterior.Sig, 2)
	transfer := math.Pow(decay, transferSpeed)

	p.ApproxPosterior = newPosterior

	wtNormOld := math.Pow(p.NormalFactor.Sig, -2)
	wtFromNormOld := transfer * wtNormOld
	wtLogistcs := 0.0
	for _, term := range p.LogisticFactors {
		wtLogistcs += term.getWeight()
	}
	wtFromTransfers := (1.0 - transfer) * (wtNormOld + wtLogistcs)
	wtTotal := wtFromNormOld + wtFromTransfers

	p.NormalFactor.Mu = (wtFromNormOld*p.NormalFactor.Mu + wtFromTransfers*p.ApproxPosterior.Mu) / wtTotal
	p.NormalFactor.Sig = math.Sqrt(1.0 / (decay * wtTotal))

	for _, r := range p.LogisticFactors {
		r.w_out *= transfer * decay
	}
}

func (p *Player) updateRatingWithLogistic(performance Rating, maxHistory uint) {
	if len(p.LogisticFactors) >= int(maxHistory) {
		logistic := p.LogisticFactors[0]
		p.LogisticFactors = p.LogisticFactors[1:]
		wn := math.Pow(p.NormalFactor.Sig, -2)
		wl := logistic.getWeight()

		p.NormalFactor.Mu = (wn*p.NormalFactor.Mu + wl*logistic.mu) / (wn - wl)
		p.NormalFactor.Sig = math.Sqrt(1.0 / (wn + wl))
	}

	p.LogisticFactors = append(p.LogisticFactors, tanhFromRating(&performance))

	newRating := p.approxPosterior(performance.Sig)
	p.updateRating(newRating, performance.Mu)
}

func (p *Player) approxPosterior(perfSig float64) Rating {
	normalWeight := math.Pow(p.NormalFactor.Sig, -2)
	logisticsCopy := make([]TanhTerm, len(p.LogisticFactors))
	copy(logisticsCopy, p.LogisticFactors) // TODO: check if this is needed
	mu := robustAverage(logisticsCopy, -p.NormalFactor.Mu*normalWeight, normalWeight)
	sig := math.Sqrt(1 / (math.Pow(p.ApproxPosterior.Sig, -2) + math.Pow(perfSig, -2)))
	return Rating{Mu: mu, Sig: sig}
}
