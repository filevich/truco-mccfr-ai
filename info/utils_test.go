package info_test

import (
	"testing"

	"github.com/filevich/truco-cfr/info"
	"github.com/truquito/truco/pdt"
)

func TestRix(t *testing.T) {
	p, _ := pdt.NuevaPartida(pdt.A20, []string{"Alice", "Anna"}, []string{"Bob", "Ben"}, 1, true)

	{
		// la mano es alice,
		// entonces:
		tests := []struct {
			who         string
			expectedRix int
		}{
			{"Alice", 0},
			{"Bob", 1},
			{"Anna", 2},
			{"Ben", 3},
		}

		for _, test := range tests {
			m := p.Manojo(test.who)
			rix := info.RIX(p, m)
			if ok := rix == test.expectedRix; !ok {
				t.Errorf("rix de %s: got:%d exp:%d; no es el esperado\n",
					m.Jugador.ID,
					rix,
					test.expectedRix)
			}
		}
	}

	{
		p.Ronda.ElMano = pdt.JIX(1)
		// la mano es bob,
		// entonces:
		tests := []struct {
			who         string
			expectedRix int
		}{
			{"Alice", 3},
			{"Bob", 0},
			{"Anna", 1},
			{"Ben", 2},
		}

		for _, test := range tests {
			m := p.Manojo(test.who)
			rix := info.RIX(p, m)
			if ok := rix == test.expectedRix; !ok {
				t.Errorf("rix de %s: got:%d exp:%d; no es el esperado\n",
					m.Jugador.ID,
					rix,
					test.expectedRix)
			}
		}
	}

	{
		p.Ronda.ElMano = pdt.JIX(2)
		// la mano es bob,
		// entonces:
		tests := []struct {
			who         string
			expectedRix int
		}{
			{"Alice", 2},
			{"Bob", 3},
			{"Anna", 0},
			{"Ben", 1},
		}

		for _, test := range tests {
			m := p.Manojo(test.who)
			rix := info.RIX(p, m)
			if ok := rix == test.expectedRix; !ok {
				t.Errorf("rix de %s: got:%d exp:%d; no es el esperado\n",
					m.Jugador.ID,
					rix,
					test.expectedRix)
			}
		}
	}

	{
		p.Ronda.ElMano = pdt.JIX(3)
		// la mano es bob,
		// entonces:
		tests := []struct {
			who         string
			expectedRix int
		}{
			{"Alice", 1},
			{"Bob", 2},
			{"Anna", 3},
			{"Ben", 0},
		}

		for _, test := range tests {
			m := p.Manojo(test.who)
			rix := info.RIX(p, m)
			if ok := rix == test.expectedRix; !ok {
				t.Errorf("rix de %s: got:%d exp:%d; no es el esperado\n",
					m.Jugador.ID,
					rix,
					test.expectedRix)
			}
		}
	}
}
