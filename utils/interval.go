// SOURCE --> "github.com/mndrix/bpci"

package utils

import "math"

// const z = 1.96
const z = 1.64

/*

1.96 is used because the 95% confidence interval has only 2.5% on each side

1% -> phi=0.995 -> z=2.57
5% -> phi=0.975 -> z=1.96
10% -> phi=0.95 -> z=1.64

*/

// Wald returns the upper and lower bound of a 95% Wald confidence
// interval where n is the number of trials and x is the number of
// successes.
//
// This method is mostly useful as a comparison.  It's not recommended
// for either large or small n.
func Wald(x int, n int) (upper float64, lower float64) {
	pHat := float64(x) / float64(n)
	delta := z * math.Sqrt(pHat*(1-pHat)/float64(n))
	upper = pHat + delta
	if upper > 1.0 {
		upper = 1.0
	}
	lower = pHat - delta
	if lower < 0.0 {
		lower = 0.0
	}
	return
}

// WaldAdjusted returns the upper and lower bound of an adjusted 95%
// Wald confidence interval ("add two successes and two failures")
// where n is the number of trials and x is the number of successes.
//
// For small n, this approximation is much more accurate than the Wald
// interval.  When x/n is near 0 or 1, this method gives more accurate
// results than the score method.
//
// This method should work well for all values of n and x/n.
//
// For justification, see "Approximate is Better than Exact for
// Interval Estimation of Binomial Proportions" by Agresti and Coull
func WaldAdjusted(x int, n int) (float64, float64) {
	return Wald(x+2, n+4)
}

// Score returns the upper and lower bound of a 95% score confidence
// interval where n is the number of trials and x is the number of
// successes.
//
// This method performs very well as along as x>1 and x<n-1.  See
// ScoreAdjusted for a version that also performs well near these
// bounds.
func Score(x int, n int) (upper float64, lower float64) {
	nf := float64(n)
	pHat := float64(x) / nf
	delta := z * math.Sqrt((pHat*(1-pHat)+z*z/(4*nf))/nf)
	upper = (pHat + (z*z)/(2*nf) + delta) / (1 + z*z/nf)
	lower = (pHat + (z*z)/(2*nf) - delta) / (1 + z*z/nf)
	return
}

// ScoreAdjusted is like Score but applies adjustments for x==1 and
// x==n-1, as suggested by Agresti and Coull.  This method should work
// well for all values of x and x/n.
func ScoreAdjusted(x int, n int) (upper float64, lower float64) {
	upper, lower = Score(x, n)
	if x == 1 {
		lower = 0.0512933 / float64(n)
	}
	if x == n-1 {
		upper = 1 - 0.0512933/float64(n)
	}
	return
}
