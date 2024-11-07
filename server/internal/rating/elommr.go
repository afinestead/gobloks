package rating

import (
	"math"
)

const SECS_PER_DAY = 86400

type ContestRatingParams struct {
	ContestWeight     float64
	PeformanceCeiling float64
}

type Standings struct {
	player *Player
	lo     uint
	hi     uint
}

type SimpleEloMMR struct {
	WeightLimit   float64
	NoobDelay     []float64
	SigLimit      float64
	DriftPerDay   float64
	SplitTies     bool
	HistoryLen    uint
	TransferSpeed float64
}

func InitSimpleEloMMR() *SimpleEloMMR {
	return &SimpleEloMMR{
		WeightLimit:   0.2,
		NoobDelay:     []float64{},
		SigLimit:      80.0,
		DriftPerDay:   0.0,
		SplitTies:     false,
		HistoryLen:    math.MaxUint,
		TransferSpeed: 1.0,
	}
}

func (mmr *SimpleEloMMR) computeWeight(contestWeight float64, n uint) float64 {
	contestWeight *= mmr.WeightLimit
	if n < uint(len(mmr.NoobDelay)) {
		contestWeight *= mmr.NoobDelay[n]
	}
	return contestWeight
}

func (mmr *SimpleEloMMR) computeSigPerf(weight float64) float64 {
	discretePerf := (1.0 + (1.0 / weight)) * mmr.SigLimit * mmr.SigLimit
	continuousPerf := mmr.DriftPerDay / weight
	return math.Sqrt(discretePerf + continuousPerf)
}

func (mmr *SimpleEloMMR) computeSigDrift(weight, deltaSec float64) float64 {
	discreteDrift := weight * mmr.SigLimit * mmr.SigLimit
	continuousDrift := mmr.DriftPerDay * deltaSec / SECS_PER_DAY
	return math.Sqrt(discreteDrift + continuousDrift)
}

func (mmr *SimpleEloMMR) IndividualUpdate(params ContestRatingParams, player *Player, muPerf float64) {
	weight := mmr.computeWeight(params.ContestWeight, player.timesPlayedExcl())
	sigPerf := mmr.computeSigPerf(weight)
	sigDrift := mmr.computeSigDrift(weight, float64(player.DeltaTime))
	player.addNoiseBest(sigDrift, mmr.TransferSpeed)
	player.updateRatingWithLogistic(Rating{Mu: muPerf, Sig: sigPerf}, mmr.HistoryLen)
}

func (mmr *SimpleEloMMR) RoundUpdate(params ContestRatingParams, standings []*Standings) {
	tanhTerms := make([]TanhTerm, len(standings))
	for ii, standing := range standings {
		weight := mmr.computeWeight(params.ContestWeight, standing.player.timesPlayedExcl())
		sigPerf := mmr.computeSigPerf(weight)
		sigDrift := mmr.computeSigDrift(weight, float64(standing.player.DeltaTime))
		standing.player.addNoiseBest(sigDrift, mmr.TransferSpeed)
		noised := standing.player.ApproxPosterior.WithNoise(sigPerf)
		tanhTerms[ii] = tanhFromRating(&noised)
	}

	var mul float64
	if mmr.SplitTies {
		mul = 1.0
	} else {
		mul = 2.0
	}

	bounds := []float64{-6000.0, 9000.0}

	for _, standing := range standings {
		f := func(x float64) (float64, float64) {
			iter := make([]evalRes, len(tanhTerms))
			for jj, term := range tanhTerms[0:standing.lo] {
				iter[jj+0] = evalLess(term, x)
			}
			for jj, term := range tanhTerms[standing.lo : standing.Hi+1] {
				iter[jj+int(standing.lo)] = evalEqual(term, x, mul)
			}
			for jj, term := range tanhTerms[standing.Hi+1:] {
				iter[jj+int(standing.Hi+1)] = evalGreater(term, x)
			}
			s := 0.0
			sp := 0.0
			for _, term := range iter {
				s += term.val
				sp += term.valPrime
			}
			return s, sp
		}
		muPerf := solveNewton(bounds, f)
		weight := mmr.computeWeight(params.ContestWeight, standing.player.timesPlayedExcl())
		sigPerf := mmr.computeSigPerf(weight)
		standing.player.updateRatingWithLogistic(Rating{Mu: muPerf, Sig: sigPerf}, mmr.HistoryLen)
	}
}
