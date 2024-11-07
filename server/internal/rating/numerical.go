package rating

import "math"

const TANH_MULTIPLIER float64 = math.Pi / 1.7320508075688772

type TanhTerm struct {
	mu    float64
	w_arg float64
	w_out float64
}

type evalRes struct {
	val, valPrime float64
}

func evalLess(term TanhTerm, x float64) evalRes {
	val, valPrime := term.baseValues(x)
	return evalRes{val - term.w_out, valPrime}
}

func evalGreater(term TanhTerm, x float64) evalRes {
	val, valPrime := term.baseValues(x)
	return evalRes{val + term.w_out, valPrime}
}

func evalEqual(term TanhTerm, x, mul float64) evalRes {
	val, valPrime := term.baseValues(x)
	return evalRes{mul * val, mul * valPrime}
}

func tanhFromRating(r *Rating) TanhTerm {
	w := TANH_MULTIPLIER / r.Sig
	return TanhTerm{
		mu:    r.Mu,
		w_arg: w * 0.5,
		w_out: w,
	}
}

func (t *TanhTerm) getWeight() float64 {
	return t.w_arg * t.w_out * 2 / math.Pow(TANH_MULTIPLIER, 2)
}

func (t *TanhTerm) baseValues(x float64) (float64, float64) {
	z := (x - t.mu) * t.w_arg
	val := -math.Tanh(z) * t.w_out
	val_prime := -math.Pow(math.Cosh(z), -2) * t.w_out * t.w_arg
	return val, val_prime
}

func robustAverage(allRatings []TanhTerm, offset, slope float64) float64 {
	bounds := []float64{-6000.0, 9000.0}

	f := func(x float64) (float64, float64) {
		allRatingsClone := make([]TanhTerm, len(allRatings))
		copy(allRatingsClone, allRatings)
		for i, term := range allRatingsClone {
			tanhZ := math.Tanh((x - term.mu) * term.w_arg)
			allRatingsClone[i] = TanhTerm{
				mu:    term.w_out * tanhZ,
				w_arg: (1.0 - tanhZ*tanhZ) * term.w_arg * term.w_out,
			}
		}
		s, sp := offset+slope*x, slope
		for _, term := range allRatingsClone {
			s += term.mu
			sp += term.w_arg
		}
		return s, sp
	}

	return solveNewton(bounds, f)
}

func solveNewton(bounds []float64, f func(float64) (float64, float64)) float64 {
	lo, hi := bounds[0], bounds[1]
	guess := 0.5 * (lo + hi)
	for {
		sum, sumPrime := f(guess)
		extrapolate := guess - sum/sumPrime
		if extrapolate < guess {
			hi = guess
			guess = math.Min(hi, math.Max(hi-0.75*(hi-lo), extrapolate))
		} else {
			lo = guess
			guess = math.Max(lo, math.Min(lo+0.75*(hi-lo), extrapolate))
		}
		if lo >= guess || guess >= hi {
			return guess
		}
	}
}
