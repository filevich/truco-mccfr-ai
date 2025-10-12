package info

import (
	"testing"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/pdt"
)

// p, _ := pdt.NuevaPartida(pdt.A20, []string{"Alice", "Anna"}, []string{"Bob", "Ben"}, 1, true)
// bs, _ := p.MarshalJSON()

var (
	verbose = true
)

func TestInfosetRondaBase(t *testing.T) {
	p, _ := pdt.Parse(`{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Ben"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":1},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"oro","valor":10},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"basto","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Anna","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Anna":2,"Ben":3,"Bob":1},"muestra":{"palo":"espada","valor":6},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]},"limiteEnvido":1}`, verbose)
	t.Log(p)
	a := abs.A1{}
	infobuilder := infosetRondaBaseFactory(a)
	i := infobuilder(p, p.Manojo("Anna"), nil)
	t.Log(i.Dump(false))
}

func eq[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestInfosetRondaLarge1NullAbs(t *testing.T) {
	p, _ := pdt.Parse(`{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Ben"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":1},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"oro","valor":10},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"basto","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Anna","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Anna":2,"Ben":3,"Bob":1},"muestra":{"palo":"espada","valor":6},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]},"limiteEnvido":1}`, verbose)
	t.Log(p)
	a := abs.Null{}

	// para Anna
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Anna"), nil)

		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 1. `muestra`
		if ok := irl.muestra == 25; !ok {
			t.Error()
		}

		// 2. `num_mano_actual`: int
		if ok := irl.numMano == 0; !ok {
			t.Error()
		}

		// 3. `rixMe` ~ RIX: who?
		if ok := irl.rixMe == 2; !ok {
			t.Error()
		}

		// 4. `turno` ~ RIX who?
		if ok := irl.rixTurno == 0; !ok {
			t.Error()
		}

		// 5. `ManojosEnJuego` quiénes se fueron al mazo y quiénes siguen en pie?
		{
			exp := []bool{true, true, true, true}
			if ok := eq(irl.manojosEnJuego, exp); !ok {
				t.Error()
			}
		}

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// cartas alice
			// (1,basto;5,espada;7,oro) --> (00;24;36) --> (2*97*157) --> 30458
			// (10,basto;7,basto;7,copa) --> (07;06;16) --> (19*17*59) --> 19057
			exp := []int{30458, 19057}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 19057; !ok {
			t.Error()
		}

		// 7. tiradas
		{
			exp := [][]int{{}, {}}
			for i, e := range exp {
				if ok := eq(irl.tiradasCartas[i], e); !ok {
					t.Error()
				}
			}
		}

		{
			exp := [][]int{{}, {}}
			for i, e := range exp {
				if ok := eq(irl.tiradasWho[i], e); !ok {
					t.Error()
				}
			}
		}

		// 8. historial
		{
			expQuien := []int{}
			expQue := []string{}
			expCuanto := []int{}
			for i := range expQue {
				if ok := expQue[i] == irl.historialQue[i]; !ok {
					t.Error()
				}
				if ok := expQuien[i] == irl.historialQuien[i]; !ok {
					t.Error()
				}
				if ok := expCuanto[i] == irl.historialCuanto[i]; !ok {
					t.Error()
				}
			}
		}

		// 9. chi
		if ok := irl.ChiLen() == 1; !ok {
			t.Error()
		}
	}

	// para Alice
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Alice"), nil)
		irl, _ := i.(*InfosetRondaXXLarge)

		// 3. `rixMe` ~ RIX: who?
		if ok := irl.rixMe == 0; !ok {
			t.Error()
		}

		// 9. chi: 3C + E|RE|FE + T + M = 8
		if ok := irl.ChiLen() == 8; !ok {
			t.Error()
		}
	}

	// para Bob
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Bob"), nil)
		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 1. `muestra`
		if ok := irl.muestra == 25; !ok {
			t.Error()
		}

		// 2. `num_mano_actual`: int
		if ok := irl.numMano == 0; !ok {
			t.Error()
		}

		// 3. `rixMe` ~ RIX: who?
		if ok := irl.rixMe == 1; !ok {
			t.Error()
		}

		// 4. `turno` ~ RIX who?
		if ok := irl.rixTurno == 0; !ok {
			t.Error()
		}

		// 5. `ManojosEnJuego` quiénes se fueron al mazo y quiénes siguen en pie?
		{
			exp := []bool{true, true, true, true}
			if ok := eq(irl.manojosEnJuego, exp); !ok {
				t.Error()
			}
		}

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// (1,esp;10,oro;10,copa) --> (20;37;17) --> (61*73*163) --> 725839
			// (6,copa;3,copa;12,copa) --> (15;12;19) --> (41*53*71) --> 154283
			exp := []int{725839, 154283}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 725839; !ok {
			t.Error()
		}

		// 7. tiradas
		{
			exp := [][]int{{}, {}}
			for i, e := range exp {
				if ok := eq(irl.tiradasCartas[i], e); !ok {
					t.Error()
				}
			}
		}

		{
			exp := [][]int{{}, {}}
			for i, e := range exp {
				if ok := eq(irl.tiradasWho[i], e); !ok {
					t.Error()
				}
			}
		}

		// 8. historial
		{
			expQuien := []int{}
			expQue := []string{}
			expCuanto := []int{}
			for i := range expQue {
				if ok := expQue[i] == irl.historialQue[i]; !ok {
					t.Error()
				}
				if ok := expQuien[i] == irl.historialQuien[i]; !ok {
					t.Error()
				}
				if ok := expCuanto[i] == irl.historialCuanto[i]; !ok {
					t.Error()
				}
			}
		}

		// 9. chi
		if ok := irl.ChiLen() == 1; !ok {
			t.Error()
		}
	}

	// para Ben
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Ben"), nil)
		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 3. `rixMe` ~ RIX: who?
		if ok := irl.rixMe == 3; !ok {
			t.Error()
		}

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// (1,esp;10,oro;10,copa) --> (20;37;17) --> (61*73*163) --> 725839
			// (6,copa;3,copa;12,copa) --> (15;12;19) --> (41*53*71) --> 154283
			exp := []int{725839, 154283}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 154283; !ok {
			t.Error()
		}

		// 9. chi: M+F
		if ok := irl.ChiLen() == 2; !ok {
			t.Error()
		}
	}
}

func TestInfosetRondaLarge1A1Abs(t *testing.T) {
	p, _ := pdt.Parse(`{"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":2,"rojo":2},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":["Ben"]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"basto","valor":1},{"palo":"oro","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":1},{"palo":"oro","valor":10},{"palo":"copa","valor":10}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"basto","valor":7},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Anna","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":6},{"palo":"copa","valor":3},{"palo":"copa","valor":12}],"tiradas":[false,false,false],"ultimaTirada":-1,"jugador":{"id":"Ben","equipo":"rojo"}}],"mixs":{"Alice":0,"Anna":2,"Ben":3,"Bob":1},"muestra":{"palo":"espada","valor":6},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]},"limiteEnvido":1}`, verbose)
	t.Log(p)
	a := abs.A1{}

	// para Anna
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Anna"), nil)
		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 1. `muestra`
		if ok := irl.muestra == 25; !ok {
			t.Error()
		}

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// cartas alice
			// (1,basto;5,espada;7,oro) --> (1;2;1) --> (3*5*3) --> 45
			// (10,basto;7,basto;7,copa) --> (0;0;0) --> (2*2*2) --> 8
			exp := []int{45, 8}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 8; !ok {
			t.Error()
		}

		// 9. chi
		if ok := irl.ChiLen() == 1; !ok {
			t.Error()
		}
	}

	// para Alice
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Alice"), nil)
		irl, _ := i.(*InfosetRondaXXLarge)

		// 9. chi: 2C + E|RE|FE + T + M = 7
		// Notar que tiene solo 2 acción de carta porque son 1 muestra + 2 matas
		if ok := irl.ChiLen() == 7; !ok {
			t.Error()
		}
	}

	// para Bob
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Bob"), nil)
		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// (1,esp;10,oro;10,copa) --> (1;0;0) --> (3*2*2) --> 12
			// (6,copa;3,copa;12,copa) --> (0;0;0) --> (2*2*2) --> 8
			exp := []int{12, 8}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 12; !ok {
			t.Error()
		}
	}

	// para Ben
	{
		infobuilder := infosetRondaXXLargeFactory(a)
		i := infobuilder(p, p.Manojo("Ben"), nil)
		// t.Log(i.Dump(true))

		irl, _ := i.(*InfosetRondaXXLarge)

		// 3. `rixMe` ~ RIX: who?
		if ok := irl.rixMe == 3; !ok {
			t.Error()
		}

		// 6. `nuestrasCartas` representa nuestras cartas.
		{
			// (1,esp;10,oro;10,copa) --> (1;0;0) --> (3*2*2) --> 12
			// (6,copa;3,copa;12,copa) --> (0;0;0) --> (2*2*2) --> 8
			exp := []int{12, 8}
			if ok := eq(irl.nuestrasCartas, exp); !ok {
				t.Error()
			}
		}

		// mi manojo
		if ok := irl._miManojoPID == 8; !ok {
			t.Error()
		}

		// 9. chi: M+F
		if ok := irl.ChiLen() == 2; !ok {
			t.Error()
		}
	}
}
