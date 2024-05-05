package dumbo

import (
	"testing"

	"github.com/filevich/truco-ai/abs"
	"github.com/truquito/gotruco/pdt"
)

func gen_partida(num_players int) *pdt.Partida {
	A := []string{"Alice", "Ariana", "Anna"}
	B := []string{"Bob", "Ben", "Bill"}
	n := num_players / 2
	p, _ := pdt.NuevaPartida(pdt.A20, A[:n], B[:n], 2, true)
	return p
}

func Test_Gen_Partida(t *testing.T) {
	num_players := 2
	p := gen_partida(num_players)
	bs, _ := p.MarshalJSON()
	t.Log(string(bs))
	t.Log(pdt.Renderizar(p))
}

func Test_dumbo_2p(t *testing.T) {
	p, _ := pdt.Parse(`{"limiteEnvido":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":1,"rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":3},{"palo":"espada","valor":2},{"palo":"basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"espada","valor":5},{"palo":"copa","valor":2},{"palo":"basto","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}}],"mixs":{"Alice":0,"Bob":1},"muestra":{"palo":"basto","valor":2},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}`, true)
	/*
	            ╔═════════════╗
	            ║             ║
	            ║             ║
	      ↓     ║    ┌2─┐     ║
	     Bob    ║    │Ba│     ║  Alice
	  ┌5─┐2─┐11┐║    └──┘     ║┌3─┐2─┐4─┐
	  │Es│Co│Ba│║             ║│Or│Es│Ba│
	  └──┘──┘──┘║             ║└──┘──┘──┘
	            ╚═════════════╝
	*/
	p.Ronda.Turno = 1

	p.Cmd("Bob 5 espada")

	alice := p.Manojo("alice")

	{
		j, _ := pdt.ParseJugada(p, "alice 4 basto")
		if ok := IsDumbo(p, alice, j, abs.A2{}); !ok {
			t.Fatal("tirar el 4 de oro deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "alice 4 basto")
		if ok := !IsDumbo(p, alice, j, abs.A1{}); !ok {
			t.Fatal("tirar el 4 de oro NO deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "alice 3 oro")
		if ok := !IsDumbo(p, alice, j, abs.A2{}); !ok {
			t.Fatal("tirar el 3 de oro NO deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "alice 3 oro")
		if ok := !IsDumbo(p, alice, j, abs.A1{}); !ok {
			t.Fatal("tirar el 3 de oro NO deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "alice 2 espada")
		if ok := !IsDumbo(p, alice, j, abs.A2{}); !ok {
			t.Fatal("tirar el 2 de espada NO deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "alice 2 espada")
		if ok := !IsDumbo(p, alice, j, abs.A1{}); !ok {
			t.Fatal("tirar el 2 de espada NO deberia considerarse dumbo")
		}
	}

}

func Test_es_ultimo(t *testing.T) {
	p, _ := pdt.Parse(`{"limiteEnvido":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":6},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"copa","valor":1},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":7},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"basto","valor":3},{"palo":"oro","valor":12}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ben","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":11},{"palo":"basto","valor":5},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Anna","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":6},{"palo":"basto","valor":1},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bill","equipo":"rojo"}}],"mixs":{"Alice":0,"Anna":4,"Ariana":2,"Ben":3,"Bill":5,"Bob":1},"muestra":{"palo":"copa","valor":2},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}`, true)
	/*
									┌11┐5─┐11┐    ┌6─┐3─┐12┐
									│Or│Ba│Co│    │Or│Ba│Or│
									└──┘──┘──┘    └──┘──┘──┘
										Anna          Ben
							╔══════════════════════════════╗
							║                              ║
							║                              ║
							║             ┌2─┐             ║
				Bill  ║             │Co│             ║  Ariana
		┌6─┐1─┐10┐║             └──┘             ║┌7─┐2─┐7─┐
		│Ba│Ba│Es│║                              ║│Or│Or│Co│
		└──┘──┘──┘║                              ║└──┘──┘──┘
							╚══════════════════════════════╝
										Alice          Bob
									┌10┐6─┐10┐    ┌11┐1─┐4─┐
									│Ba│Es│Or│    │Ba│Co│Es│
									└──┘──┘──┘    └──┘──┘──┘
	*/

	p.Cmd("Alice 10 basto")
	p.Cmd("Bob 11 basto")
	p.Cmd("Ariana mazo")
	p.Cmd("Ben mazo")
	p.Cmd("Anna 11 oro")

	// ahora Bill es el ultimo en tirar?
	// era el ultimo en tirar de esta mano?
	bill := p.Manojo("Bill")
	bob := p.Manojo("Bob")
	alice := p.Manojo("Alice")
	anna := p.Manojo("Anna")

	if ok := esElUltimoEnTirar(p, bill); !ok {
		t.Fatal("Deberia ser el ultimo en tirar")
	}

	p.Cmd("Bill 6 basto")
	t.Log(pdt.Renderizar(p))

	// gana bob
	// recordar: ben y ariana se fueron

	if ok := !esElUltimoEnTirar(p, bob); !ok {
		t.Fatal("NO Deberia ser el ultimo en tirar")
	}

	if ok := !esElUltimoEnTirar(p, anna); !ok {
		t.Fatal("NO Deberia ser el ultimo en tirar")
	}

	if ok := !esElUltimoEnTirar(p, bill); !ok {
		t.Fatal("NO Deberia ser el ultimo en tirar")
	}

	p.Cmd("Bob 1 copa")
	p.Cmd("Anna 5 basto")

	if ok := !esElUltimoEnTirar(p, bill); !ok {
		t.Fatal("NO Deberia ser el ultimo en tirar")
	}

	p.Cmd("Bill 1 basto")

	if ok := esElUltimoEnTirar(p, alice); !ok {
		t.Fatal("*Deberia ser el ultimo en tirar")
	}

	t.Log(pdt.Renderizar(p))

}

func Test_dumbo_6p(t *testing.T) {
	p, _ := pdt.Parse(`{"limiteEnvido":4,"puntuacion":20,"puntajes":{"azul":0,"rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"azul":3,"rojo":3},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noGritadoAun"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":10},{"palo":"espada","valor":6},{"palo":"oro","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":11},{"palo":"copa","valor":1},{"palo":"espada","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":5},{"palo":"oro","valor":2},{"palo":"copa","valor":7}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ariana","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"oro","valor":6},{"palo":"basto","valor":3},{"palo":"oro","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Ben","equipo":"rojo"}},{"seFueAlMazo":false,"cartas":[{"palo":"basto","valor":12},{"palo":"oro","valor":7},{"palo":"copa","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Anna","equipo":"azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"copa","valor":4},{"palo":"basto","valor":1},{"palo":"espada","valor":10}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bill","equipo":"rojo"}}],"mixs":{"Alice":0,"Anna":4,"Ariana":2,"Ben":3,"Bill":5,"Bob":1},"muestra":{"palo":"copa","valor":2},"manos":[{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]},{"resultado":"indeterminado","ganador":"","cartasTiradas":[]}]}}`, true)

	p.Cmd("Alice 10 basto")
	p.Cmd("Bob 11 basto")
	p.Cmd("Ariana mazo")
	p.Cmd("Ben mazo")
	p.Cmd("Anna 12 basto")

	t.Log(pdt.Renderizar(p))
	t.Log("max tiradas:", _maxTiradasStr(p, abs.A1{}))

	// [Rojo] - Bob - 11 de Basto - a1 ~> bucket #0
	// [Azul] - Anna - 12 de Oro - a1 ~> bucket #0
	// va ganando azul

	// ahora bill para ganar puede: tirar el 1 de basto o el 4 de la muestra
	// tirar el 4 de la muestra seria dummy

	bill := p.Manojo("Bill")

	{
		j, _ := pdt.ParseJugada(p, "Bill 4 copa")
		if ok := IsDumbo(p, bill, j, abs.A1{}); !ok {
			t.Fatal("tirar el 4 de copa deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "Bill 1 basto")
		if ok := !IsDumbo(p, bill, j, abs.A1{}); !ok {
			t.Fatal("tirar el 1 de basto NO deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "Bill 10 espada")
		if ok := !IsDumbo(p, bill, j, abs.A1{}); !ok {
			t.Fatal("tirar el 10 de espada NO deberia considerarse dumbo")
		}
	}

	p.Cmd("Bill 1 basto")

	// segunda mano
	p.Cmd("Bill 10 espada")
	p.Cmd("Alice 10 oro")
	p.Cmd("Bob 1 copa")

	t.Log(pdt.Renderizar(p))
	t.Log("max tiradas:", _maxTiradasStr(p, abs.A2{}))

	{
		muestra := p.Ronda.Muestra
		c1 := &pdt.Carta{Palo: pdt.Copa, Valor: 11}
		c2 := &pdt.Carta{Palo: pdt.Oro, Valor: 7}
		t.Log("buck abs 11 Co ~>", abs.A2{}.Abstract(c1, &muestra))
		t.Log("buck abs 7 Or ~>", abs.A2{}.Abstract(c2, &muestra))
	}

	anna := p.Manojo("anna")

	{
		j, _ := pdt.ParseJugada(p, "Anna 11 copa")
		if ok := IsDumbo(p, anna, j, abs.A2{}); !ok {
			t.Fatal("tirar el 11 de copa deberia considerarse dumbo")
		}
	}

	{
		j, _ := pdt.ParseJugada(p, "Anna 7 oro")
		if ok := !IsDumbo(p, anna, j, abs.A2{}); !ok {
			t.Fatal("tirar el 7 de oro NO deberia considerarse dumbo")
		}
	}

}
