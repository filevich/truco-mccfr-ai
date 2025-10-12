package abs_test

import (
	"testing"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/pdt"
)

func TestAbstraccionNull(t *testing.T) {
	var (
		abs     abs.IAbstraction = abs.A3{}
		carta   pdt.Carta        = pdt.NuevaCarta(pdt.CartaID(13))
		muestra *pdt.Carta       = &carta
	)

	for i := 0; i < 40; i++ {
		c := pdt.NuevaCarta(pdt.CartaID(i))
		// exp := i
		got := abs.Abstract(&c, muestra)
		t.Logf("i:%d muestra: `%s` carta:`%s` abs_a3:%d", i, muestra, c, got)
		// if ok := got == exp; !ok {
		// 	t.Errorf("el id no es el esperado. got:%d exp:%d", got, exp)
		// }
	}

}
