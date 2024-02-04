package main

import (
	"encoding/json"
	"os"

	"github.com/filevich/combinatronics"
	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/truco/pdt"
)

func inc(dict map[int]int, key int) {
	if _, ok := dict[key]; !ok {
		dict[key] = 0
	}
	dict[key] += 1
}

func calcSumPower(m *pdt.Manojo, muestra pdt.Carta) int {
	total := 0
	for _, c := range m.Cartas {
		total += c.CalcPoder(muestra)
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
