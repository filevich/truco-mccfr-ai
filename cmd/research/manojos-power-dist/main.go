package main

import (
	"encoding/json"
	"os"

	"github.com/filevich/combinatronics"
	"github.com/filevich/truco-mccfr-ai/utils"
	"github.com/truquito/gotruco/pdt"
)

func calcPoder(c *pdt.Carta, muestra pdt.Carta) int {
	var poder int

	if c.EsPieza(muestra) {
		switch c.Valor {
		case 2:
			poder = 18
		case 4:
			poder = 17
		case 5:
			poder = 16
		case 11:
			poder = 15
		case 10:
			poder = 14
		case 12:
			valeComo := &pdt.Carta{Palo: c.Palo, Valor: muestra.Valor}
			poder = calcPoder(valeComo, muestra)
		}

	} else if c.Palo == pdt.Espada && c.Valor == 1 {
		poder = 13
	} else if c.Palo == pdt.Basto && c.Valor == 1 {
		poder = 12
	} else if c.Palo == pdt.Espada && c.Valor == 7 {
		poder = 11
	} else if c.Palo == pdt.Oro && c.Valor == 7 {
		poder = 10
		// Chicas
	} else if c.Valor == 3 {
		poder = 9
	} else if c.Valor == 2 {
		poder = 8
	} else if c.Valor == 1 {
		poder = 7
	} else if c.Valor == 12 {
		poder = 6
	} else if c.Valor == 11 {
		poder = 5
	} else if c.Valor == 10 {
		poder = 4
	} else if c.Valor == 7 {
		poder = 3
	} else if c.Valor == 6 {
		poder = 2
	} else if c.Valor == 5 {
		poder = 1
	} else if c.Valor == 4 {
		poder = 0
	}

	return poder
}

func inc(dict map[int]int, key int) {
	if _, ok := dict[key]; !ok {
		dict[key] = 0
	}
	dict[key] += 1
}

func calcSumPower(m *pdt.Manojo, muestra pdt.Carta) int {
	total := 0
	for _, c := range m.Cartas {
		total += calcPoder(c, muestra)
	}
	return total
}

func save(dict map[int]int, filename string) {
	file, _ := json.MarshalIndent(dict, "", " ")
	_ = os.WriteFile(filename, file, 0644)
}

func main() {

	n := 40 // full deck
	// n := 12 // mini-truco

	ids := make([]int, n)
	for i := range ids {
		ids[i] = i
	}

	distEnvido := make(map[int]int)   // dist:Poder -> count
	distFlores := make(map[int]int)   // dist:Poder -> count
	distPowerSum := make(map[int]int) // dist:Poder -> count

	p, _ := pdt.NuevaPartida(pdt.A20, []string{"Alice"}, []string{"Bob"}, 4, false)
	m := p.Ronda.Manojo("Alice")

	// todas las muestras posibles
	for _, muestraID := range ids {
		resto := utils.CopyWithoutThese(ids, muestraID)
		// todos mis manojos posibles
		for _, miManojoIDs := range combinatronics.Combs(resto, 3) {
			c0 := pdt.NuevaCarta(pdt.CartaID(miManojoIDs[0]))
			c1 := pdt.NuevaCarta(pdt.CartaID(miManojoIDs[1]))
			c2 := pdt.NuevaCarta(pdt.CartaID(miManojoIDs[2]))
			m.Cartas = [3]*pdt.Carta{&c0, &c1, &c2}
			muestra := pdt.NuevaCarta(pdt.CartaID(muestraID))
			p.Ronda.SetMuestra(muestra)

			tieneFlor, _ := m.TieneFlor(muestra)
			if tieneFlor {
				pts, _ := m.CalcFlor(muestra)
				inc(distFlores, pts)
			} else {
				pts := m.CalcularEnvido(muestra)
				inc(distEnvido, pts)
			}

			pts := calcSumPower(m, muestra)
			inc(distPowerSum, pts)
		}
	}

	save(distEnvido, "/tmp/dist-envido.json")
	save(distFlores, "/tmp/dist-flor.json")
	save(distPowerSum, "/tmp/dist-power-sum.json")
}
