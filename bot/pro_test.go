package bot

import "testing"

func TestCDF(t *testing.T) {
	pro := Pro{}
	pro.Initialize()
	if ok := .8 < pro.powerDist.cdf(84); !ok {
		t.Error("cdf for value 84 is lower than expected")
	}
	if ok := pro.powerDist.cdf(32) < .1; !ok {
		t.Error("cdf for value 32 is higher than expected")
	}
	{
		p := pro.powerDist.cdf(50)
		if ok := .45 < p && p < .55; !ok {
			t.Error("unexpected cdf for value 50")
		}
	}
	if ok := 0 < pro.powerDist.cdf(30); !ok {
		t.Error("cdf for value 30 should be positive")
	}
	if ok := pro.powerDist.cdf(99) == 1; !ok {
		t.Error("cdf for value 99 should be 1")
	}

	// for i := 30; i < 99+1; i++ {
	// 	t.Logf("CDF(%d)=%.3f\n", i, pro.powerDist.cdf(i))
	// }
}

func TestProbDare(t *testing.T) {
	pro := Pro{}
	pro.Initialize()
	lowerBoundDare := float32(.1)
	alpha := float32(3)
	beta := float32(0)
	gamma := float32(2)
	k := float32(0.05)
	for i := 30; i < 99+1; i++ {
		t.Logf("for x=%d -> CDF(%d)=%.3f -> probDareLineal(%d,%.2f)=%.2f & probDareTanh(%d,%.0f,%.0f,%.2f)=%.2f\n",
			i,
			i, pro.powerDist.cdf(i),
			i, lowerBoundDare, pro.powerDist.probDareLineal(i, lowerBoundDare),
			i, alpha, beta, k, pro.powerDist.probDareTanh(i, alpha, beta, gamma, k),
		)
	}
}
