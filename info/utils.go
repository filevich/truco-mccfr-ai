package info

import (
	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

// Nota: `maxWe` y `maxOp` retorna la abstraccion!!
// O sea, esta funcion es incompatible con la abstracción `Null` debido a que
// en la abstraccion `Null` el bucket no está relacionado con el poder
// de la carta; (como sí sucede con A1, A2 y A3)
// ---------------------------
// retorna true si Vamos ganando, desde la perspectiva de `manojo`
// para decidir si vamos ganando|perdiendo|empatados NO USA ABSTRACCION
// pero Max_us, Max_op si está abstraido
// ---------------------------
// el `vamos` NO depende de la abstracción
func vamos(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	abs abs.IAbstraccion,

) (maxWeAbstracto, maxOpAbstracto int, vamosExacto string) {

	// retornos
	maxWeAbstracto = -1
	maxOpAbstracto = -1
	vamosExacto = ""

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

	// seteo los valores máximos
	// estos valores tienen como dominio el rango de los buckets de la abstraccion
	if maxCarta[weEquipo] != nil {
		maxWeAbstracto = abs.Abstraer(&maxCarta[weEquipo].Carta, &p.Ronda.Muestra)
	}
	if maxCarta[opEquipo] != nil {
		maxOpAbstracto = abs.Abstraer(&maxCarta[opEquipo].Carta, &p.Ronda.Muestra)
	}

	// finalmente comparo los PODERES exactos y determino si vamos ganando o no
	if maxCarta[weEquipo] == nil || maxCarta[opEquipo] == nil {
		vamosExacto = "?"
	} else {
		if maxPoder[weEquipo] < maxPoder[opEquipo] {
			vamosExacto = "perdiendo"
		} else if maxPoder[weEquipo] == maxPoder[opEquipo] {
			vamosExacto = "empatados"
		} else {
			vamosExacto = "ganando"
		}
	}

	return maxWeAbstracto, maxOpAbstracto, vamosExacto
}

// índice *relativo AL MANO* de un jugador dado
func RIX(p *pdt.Partida, m *pdt.Manojo) int {
	n := len(p.Ronda.Manojos)
	manoIx := int(p.Ronda.ElMano)
	manojoIx := p.Ronda.MIXS[m.Jugador.ID]
	return utils.Mod(manojoIx-manoIx, n)
}
