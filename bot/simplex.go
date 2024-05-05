package bot

import (
	"math"

	"github.com/filevich/truco-ai/info"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type SimpleX struct {
	inGameID string
}

func (b *SimpleX) Initialize() {}

func (b *SimpleX) Free() {}

func (b *SimpleX) UID() string {
	return "SimpleX"
}

func (b *SimpleX) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *SimpleX) ResetCatch() {}

// retorna la carta mas chicas y mas poderosa en su poder SEGUN b.Abs.Abstraer(...) !
// si no tiene getMinMaxCartas en su poder, retorna -1,-1
func (b *SimpleX) getMinMaxCartas(p *pdt.Partida) (min, max *pdt.Carta) {
	minPoder, maxPoder := math.MaxInt32, math.MinInt32
	minCartaIx, maxCartaIx := -1, -1
	m := p.Manojo(b.inGameID)
	for cix, tirada := range m.Tiradas {
		if !tirada {
			poder := m.Cartas[cix].CalcPoder(p.Ronda.Muestra) // antes
			// poder := b.Abs.Abstraer(m.Cartas[cix], &p.Ronda.Muestra) // ahora
			if poder > maxPoder {
				maxPoder = poder
				maxCartaIx = cix
			}
			if poder < minPoder {
				minPoder = poder
				minCartaIx = cix
			}
		}
	}

	return m.Cartas[minCartaIx], m.Cartas[maxCartaIx]
}

func (b *SimpleX) jugarCarta(p *pdt.Partida) pdt.IJugada {
	m := p.Manojo(b.inGameID)
	_, maxOpCarta, vamos := info.Vamos(p, m)

	switch vamos {
	case "ganando":
		// tiro la mas chica que tenga porque no hay necesidad
		minC, _ := b.getMinMaxCartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *minC,
		}

	// es indiferente la abstraccion porque:
	// el caso ganando/perdiendo/empatados se calcula segun Carta.CalcPoder(...)
	// la cual es una funcion independiente de la abstraccion

	case "perdiendo", "empatados", "?":

		// la idea es tirar la MINIMA carta TIRABLE en nuestro poder en esta mano
		// si la tengo yo -> entonces la tiro
		// si no la tengo yo -> entonces tiro la de menor poder

		// empezando por mi, me fijo de mi equipo quién tiene la carta más baja
		// que le gane a la `maxOpCarta`

		minPoderCartaSuperior := math.MaxInt
		esMia := false
		cartaIdx := -1

		for i := p.Ronda.MIXS[m.Jugador.ID]; i < len(p.Ronda.Manojos); i++ {
			manojo := p.Ronda.Manojos[i]
			if esDeMiEquipo := manojo.Jugador.Equipo == m.Jugador.Equipo; !esDeMiEquipo {
				continue
			}

			for cix, c := range manojo.Cartas {
				if manojo.Tiradas[cix] {
					continue
				}
				// no ha sido tirada.
				// es superior?
				poder := c.CalcPoder(p.Ronda.Muestra)
				maxOp := -1
				if maxOpCarta != nil {
					maxOp = maxOpCarta.CalcPoder(p.Ronda.Muestra)
				}
				if esSuperior := poder > maxOp; esSuperior {
					// ok, es menor de la mejor encontrada hasta ahora?
					if esMasEconomica := poder < minPoderCartaSuperior; esMasEconomica {
						minPoderCartaSuperior = poder
						cartaIdx = cix
						esMia = manojo.Jugador.ID == m.Jugador.ID
					}
				}
			}
		}

		if esMia {
			return &pdt.TirarCarta{
				JID:   m.Jugador.ID,
				Carta: *m.Cartas[cartaIdx],
			}
		}

		minC, _ := b.getMinMaxCartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *minC,
		}

	//
	default:
		return nil
	}
}

func (b *SimpleX) jugarLaFlor(p *pdt.Partida) pdt.IJugada {
	// si no cante -> la canto
	// si tengo que responder a una apuesta: respondo quiero solo cuando mi flor
	// es > a la flor media (24 = 47/2 (max flor 47 = 30+29+28))
	canteMiFlor := true
	for _, jid := range p.Ronda.Envite.SinCantar {
		if jid == b.inGameID {
			canteMiFlor = false
		}
	}

	if !canteMiFlor {
		// ojo que si soy el ultimo, podria considerar aumentar la apuesta aqui/ahora
		return &pdt.CantarFlor{
			JID: b.inGameID,
		}
	}

	// si no -> tengo que responder por quiero/no-quiero/aumentar-la-apuesta
	poder, _ := p.Manojo(b.inGameID).CalcFlor(p.Ronda.Muestra)
	if poder > 24 {
		return &pdt.ResponderQuiero{
			JID: b.inGameID,
		}
	}

	return &pdt.ResponderNoQuiero{
		JID: b.inGameID,
	}
}

func (b *SimpleX) jugarElEnvido(p *pdt.Partida) pdt.IJugada {
	// envido maximo: 30 + 7 = 37
	// responso si solo cuando: mi envido es mayor a 37/2 = 18
	e := p.Manojo(b.inGameID).CalcularEnvido(p.Ronda.Muestra)
	if e > 18 {
		// debo considerar aumenar o no la apuesta aqui
		return &pdt.ResponderQuiero{
			JID: b.inGameID,
		}
	}

	return &pdt.ResponderNoQuiero{
		JID: b.inGameID,
	}
}

func (b *SimpleX) jugarElTruco(p *pdt.Partida) pdt.IJugada {

	cantMuestras := 0
	for _, c := range p.Manojo(b.inGameID).Cartas {
		if c.EsPieza(p.Ronda.Muestra) {
			cantMuestras++
		}
	}

	if nadieGritoNada := p.Ronda.Truco.Estado == pdt.NOGRITADOAUN; nadieGritoNada {
		// grito?
		if condicionGritoTruco := cantMuestras >= 1; condicionGritoTruco {
			return &pdt.GritarTruco{JID: b.inGameID}
		}
	} else if nosPropusieronTruco := p.Ronda.Truco.Estado == pdt.TRUCO &&
		p.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(b.inGameID).Jugador.Equipo; nosPropusieronTruco {
		// acepto y me quedo
		// acepto y subo
		// no-acepto
		if noAcepto := cantMuestras == 0; noAcepto {
			return &pdt.ResponderNoQuiero{JID: b.inGameID}
		} else if aceptoYMeQuedo := cantMuestras == 1; aceptoYMeQuedo {
			return &pdt.ResponderQuiero{JID: b.inGameID}
		} else if aceptoYSubo := cantMuestras > 1; aceptoYSubo {
			return &pdt.GritarReTruco{JID: b.inGameID}
		}
	} else if nosPropusieronReTruco := p.Ronda.Truco.Estado == pdt.RETRUCO &&
		p.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(b.inGameID).Jugador.Equipo; nosPropusieronReTruco {
		// acepto y me quedo
		// acepto y subo
		// no-acepto
		if noAcepto := cantMuestras < 2; noAcepto {
			return &pdt.ResponderNoQuiero{JID: b.inGameID}
		} else if aceptoYMeQuedo := cantMuestras == 2; aceptoYMeQuedo {
			return &pdt.ResponderQuiero{JID: b.inGameID}
		} else if aceptoYSubo := cantMuestras > 2; aceptoYSubo {
			return &pdt.GritarVale4{JID: b.inGameID}
		}
	} else if nosPropusieronVale4 := p.Ronda.Truco.Estado == pdt.VALE4 &&
		p.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo != p.Manojo(b.inGameID).Jugador.Equipo; nosPropusieronVale4 {
		// acepto y me quedo
		// acepto y subo
		// no-acepto
		if noAcepto := cantMuestras < 3; noAcepto {
			return &pdt.ResponderNoQuiero{JID: b.inGameID}
		} else if aceptoYMeQuedo := cantMuestras == 3; aceptoYMeQuedo {
			return &pdt.ResponderQuiero{JID: b.inGameID}
		}
	}

	return nil
}

// pre: el jugador no se fue al mazo
func (b *SimpleX) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	b.inGameID = inGameID

	if pdt.LaFlorEsRespondible(p) {
		if _, err := p.Manojo(b.inGameID).CalcFlor(p.Ronda.Muestra); err == nil {
			return b.jugarLaFlor(p), 1
		}
	}

	if pdt.ElEnvidoEsRespondible(p) {
		return b.jugarElEnvido(p), 1
	}

	// debo considerar irme al mazo

	if j := b.jugarElTruco(p); j != nil {
		return j, 1
	}

	// si no: elTrucoRespondible(p)
	return b.jugarCarta(p), 1
}
