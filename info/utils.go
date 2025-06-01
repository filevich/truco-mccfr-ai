package info

import (
	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/filevich/truco-mccfr-ai/utils"
	"github.com/truquito/gotruco/pdt"
)

// Nota: `maxWe` y `maxOp` retorna la abstraccion!!
// O sea, esta funcion es incompatible con la abstracción `Null` debido a que
// en la abstraccion `Null` el bucket no está relacionado con el poder
// de la carta; (como sí sucede con A1, A2 y A3)
// ---------------------------
// retorna true si Vamos ganando, desde la perspectiva de `manojo`
// para decidir si Vamos ganando|perdiendo|empatados NO USA ABSTRACCION
// pero Max_us, Max_op si está abstraido
// ---------------------------
// el `Vamos` NO depende de la abstracción
func Vamos(

	p *pdt.Partida,
	manojo *pdt.Manojo,

) (ourMax, opMax *pdt.Carta, vamosEnum string) {

	// retornos
	vamosEnum = ""

	// determino el qué carta y qué poder tienen las cartas más poderosas tiradas
	// por mi equipo y el equipo de op.
	maxPoder := map[pdt.Equipo]int{pdt.Rojo: -1, pdt.Azul: -1}
	maxCarta := map[pdt.Equipo]*pdt.CartaTirada{pdt.Rojo: nil, pdt.Azul: nil}
	tiradas := p.Ronda.GetManoActual().CartasTiradas
	for i, tirada := range tiradas {
		poder := tirada.Carta.CalcPoder(p.Ronda.Muestra)
		equipo := p.Ronda.Manojo(tirada.Jugador).Jugador.Equipo
		if poder > maxPoder[equipo] {
			maxPoder[equipo] = poder
			maxCarta[equipo] = &tiradas[i]
		}
	}

	weEquipo := manojo.Jugador.Equipo
	opEquipo := manojo.Jugador.GetEquipoContrario()

	// finalmente comparo los PODERES exactos y determino si vamos ganando o no
	if maxCarta[weEquipo] == nil || maxCarta[opEquipo] == nil {
		vamosEnum = "?"
	} else {
		if maxPoder[weEquipo] < maxPoder[opEquipo] {
			vamosEnum = "perdiendo"
		} else if maxPoder[weEquipo] == maxPoder[opEquipo] {
			vamosEnum = "empatados"
		} else {
			vamosEnum = "ganando"
		}
	}

	var cWe, cOp *pdt.Carta = nil, nil

	if maxCarta[weEquipo] != nil {
		cWe = &maxCarta[weEquipo].Carta
	}

	if maxCarta[opEquipo] != nil {
		cOp = &maxCarta[opEquipo].Carta
	}

	return cWe, cOp, vamosEnum
}

// índice *relativo AL MANO* de un jugador dado
func RIX(p *pdt.Partida, m *pdt.Manojo) int {
	n := len(p.Ronda.Manojos)
	manoIx := int(p.Ronda.ElMano)
	manojoIx := p.Ronda.MIXS[m.Jugador.ID]
	return utils.Mod(manojoIx-manoIx, n)
}

func PrimifyManojo(

	m *pdt.Manojo,
	muestra *pdt.Carta,
	a abs.IAbstraction,

) int {

	manojoPrimeID := 1
	for _, c := range m.Cartas {
		cartaAbstraida := a.Abstract(c, muestra)
		manojoPrimeID *= utils.AllPrimes[cartaAbstraida]
	}
	return manojoPrimeID

}
