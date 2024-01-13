package trucocfr_test

import (
	"testing"

	trucocfr "github.com/filevich/truco-cfr"
	"github.com/truquito/truco/pdt"
)

func TestAbstraccionZero(t *testing.T) {
	var (
		abs     trucocfr.IAbstraccion = trucocfr.Zero{}
		muestra *pdt.Carta            = nil
	)

	for i := 0; i < 40; i++ {
		c := pdt.NuevaCarta(pdt.CartaID(i))
		exp := i + 1
		got := abs.Abstraer(&c, muestra)
		t.Logf("i:%d carta:%s abs_zero:%d", i, c, got)
		if ok := got == exp; !ok {
			t.Errorf("el id no es el esperado. got:%d exp:%d", got, exp)
		}
	}

}
