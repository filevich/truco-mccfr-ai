package abs_test

import (
	"testing"

	"github.com/filevich/truco-cfr/abs"
	"github.com/truquito/truco/pdt"
)

func TestAbstraccionNull(t *testing.T) {
	var (
		abs     abs.IAbstraction = abs.Null{}
		muestra *pdt.Carta       = nil
	)

	for i := 0; i < 40; i++ {
		c := pdt.NuevaCarta(pdt.CartaID(i))
		exp := i
		got := abs.Abstract(&c, muestra)
		t.Logf("i:%d carta:`%s` abs_null:%d", i, c, got)
		if ok := got == exp; !ok {
			t.Errorf("el id no es el esperado. got:%d exp:%d", got, exp)
		}
	}

}
