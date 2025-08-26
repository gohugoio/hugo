// This is just a temporary fork of https://github.com/josharian/vitter (ISC License, https://github.com/josharian/vitter/blob/main/LICENSE)
//
// This file will be removed once https://github.com/josharian/vitter/issues/1 is resolved.

package collections

import (
	"math"
	"math/rand/v2"
)

// https://getkerf.wordpress.com/2016/03/30/the-best-algorithm-no-one-knows-about/

// Copyright Kevin Lawler, released under ISC License

// _d generates an in-order uniform random sample of size 'want' from the range [0, max) using the provided PRNG.
//
// Parameters:
//   - prng: random number generator
//   - want: number of samples to select
//   - max: upper bound of the range [0, max) from which to sample
//   - choose: callback function invoked with each selected index in ascending order
//
// If the parameters are invalid (want < 0 or want > max), no samples are selected.
//
// Vitter, J.S. - An Efficient Algorithm for Sequential Random Sampling - ACM Trans. Math. Software 11 (1985), 37-57.
func _d(prng *rand.Rand, want, max int, choose func(n int)) {
	if want <= 0 || want > max {
		return
	}
	// POTENTIAL_OPTIMIZATION_POINT: Christian Neukirchen points out we can replace exp(log(x)*y) by pow(x,y)
	// POTENTIAL_OPTIMIZATION_POINT: Vitter paper points out an exponentially distributed random var can provide speed ups
	// 'a' is space allocated for the hand
	// 'n' is the size of the hand
	// 'N' is the upper bound on the random card values
	j := -1
	qu1 := -want + 1 + max
	const negalphainv = -13 // threshold parameter from Vitter's paper for algorithm selection
	threshold := -negalphainv * want

	wantf := float64(want)
	maxf := float64(max)
	ninv := 1.0 / wantf
	var nmin1inv float64
	Vprime := math.Exp(math.Log(prng.Float64()) * ninv)

	qu1real := -wantf + 1.0 + maxf
	var U, X, y1, y2, top, bottom, negSreal float64

	for want > 1 && threshold < max {
		var S int

		nmin1inv = 1.0 / (-1.0 + wantf)

		for {
			for {
				X = maxf * (-Vprime + 1.0)
				S = int(math.Floor(X))

				if S < qu1 {
					break
				}

				Vprime = math.Exp(math.Log(prng.Float64()) * ninv)
			}

			U = prng.Float64()
			negSreal = float64(-S)
			y1 = math.Exp(math.Log(U*maxf/qu1real) * nmin1inv)
			Vprime = y1 * (-X/maxf + 1.0) * (qu1real / (negSreal + qu1real))

			if Vprime <= 1.0 {
				break
			}

			y2 = 1.0
			top = -1.0 + maxf
			var limit int

			if -1+want > S {
				bottom = -wantf + maxf
				limit = -S + max
			} else {
				bottom = -1.0 + negSreal + maxf
				limit = qu1
			}

			for t := max - 1; t >= limit; t-- {
				y2 = (y2 * top) / bottom
				top--
				bottom--
			}

			if maxf/(-X+maxf) >= y1*math.Exp(math.Log(y2)*nmin1inv) {
				Vprime = math.Exp(math.Log(prng.Float64()) * nmin1inv)
				break
			}

			Vprime = math.Exp(math.Log(prng.Float64()) * ninv)
		}

		j += S + 1

		choose(j)

		max = -S + (-1 + max)
		maxf = negSreal + (-1.0 + maxf)
		want--
		wantf--
		ninv = nmin1inv

		qu1 = -S + qu1
		qu1real = negSreal + qu1real

		threshold += negalphainv
	}

	if want > 1 {
		methodA(prng, want, max, j, choose) // if i>0 then n has been decremented
	} else {
		S := int(math.Floor(float64(max) * Vprime))

		j += S + 1

		choose(j)
	}
}

// methodA is the simpler fallback algorithm used when Algorithm D's optimizations are not beneficial.
func methodA(prng *rand.Rand, want, max int, j int, choose func(n int)) {
	for want >= 2 {
		j++
		V := prng.Float64()
		quot := float64(max-want) / float64(max)
		for quot > V {
			j++
			max--
			quot *= float64(max - want)
			quot /= float64(max)
		}
		choose(j)
		max--
		want--
	}

	S := int(math.Floor(float64(max) * prng.Float64()))
	j += S + 1
	choose(j)
}
