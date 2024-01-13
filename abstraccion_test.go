package cfr_test

import (
	"testing"

	cfr "github.com/filevich/truco-cfr"
	"github.com/truquito/truco/pdt"
)

func TestAbstraccionZero(t *testing.T) {
	var (
		abs     cfr.IAbstraccion = cfr.Zero{}
		muestra *pdt.Carta       = nil
	)

	for i := 0; i < 40; i++ {
		c := pdt.NuevaCarta(pdt.CartaID(i))
		exp := i
		got := abs.Abstraer(&c, muestra)
		t.Logf("i:%d carta:%s abs_zero:%d", i, c, got)
		if ok := got == exp; !ok {
			t.Errorf("el id no es el esperado. got:%d exp:%d", got, exp)
		}
	}

}
