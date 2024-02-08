package bot

import (
	"math"

	"github.com/filevich/truco-ai/info"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type Simple struct {
	inGameID string
}

func (b *Simple) Initialize() {}

func (b *Simple) Free() {}

func (b *Simple) UID() string {
	return "Simple"
}

func (b *Simple) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *Simple) ResetCatch() {}

// retorna la carta mas chicas y mas poderosa en su poder SEGUN b.Abs.Abstraer(...) !
// si no tiene cartas en su poder, retorna -1,-1
func (b *Simple) cartas(p *pdt.Partida) (min, max *pdt.Carta) {
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

func (b *Simple) jugar_carta(p *pdt.Partida) *pdt.TirarCarta {
	m := p.Manojo(b.inGameID)
	_, maxOpCarta, vamos := info.Vamos(p, m)

	switch vamos {
	case "ganando":
		// tiro la mas chica que tenga porque no hay necesidad
		minC, _ := b.cartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *minC,
		}

	// es indiferente la abstraccion porque:
	// el caso ganando/perdiendo/empatados se calcula segun Carta.CalcPoder(...)
	// la cual es una funcion independiente de la abstraccion

	case "perdiendo", "empatados":
		// la carta mas alta mia, supera a la de op?
		// en caso afirmativo, la tiro
		// en caso negativo, no vale la pena
		minC, maxC := b.cartas(p)
		if maxC.CalcPoder(p.Ronda.Muestra) > maxOpCarta.CalcPoder(p.Ronda.Muestra) {
			return &pdt.TirarCarta{
				JID:   m.Jugador.ID,
				Carta: *maxC,
			}
		} else {
			return &pdt.TirarCarta{
				JID:   m.Jugador.ID,
				Carta: *minC,
			}
		}

	// case "?"
	default:
		// tiro yo primero, entonce tiro la mas alta
		_, maxC := b.cartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *maxC,
		}
	}

}

func (b *Simple) jugar_flor(p *pdt.Partida) pdt.IJugada {
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
			// Manojo: p.Manojo(b.inGameID),
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

func (b *Simple) jugar_envido(p *pdt.Partida) pdt.IJugada {
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

func (b *Simple) jugar_truco(p *pdt.Partida) pdt.IJugada {
	cantMuestras := 0
	for _, c := range p.Manojo(b.inGameID).Cartas {
		if c.EsPieza(p.Ronda.Muestra) {
			cantMuestras++
		}
	}

	quiero := &pdt.ResponderQuiero{
		JID: b.inGameID,
	}
	no_quiero := &pdt.ResponderNoQuiero{
		JID: b.inGameID,
	}

	condicion_quiero := cantMuestras == 1 && p.Ronda.Truco.Estado == pdt.TRUCO ||
		cantMuestras == 2 && p.Ronda.Truco.Estado == pdt.RETRUCO ||
		cantMuestras == 3 && p.Ronda.Truco.Estado == pdt.VALE4

	if condicion_quiero {
		return quiero
	}

	return no_quiero
}

// pre: el jugador no se fue al mazo
func (b *Simple) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	b.inGameID = inGameID

	if pdt.LaFlorEsRespondible(p) {
		if _, err := p.Manojo(b.inGameID).CalcFlor(p.Ronda.Muestra); err == nil {
			return b.jugar_flor(p), 1
		}
	}

	if pdt.ElEnvidoEsRespondible(p) {
		return b.jugar_envido(p), 1
	}

	// debo considerar irme al mazo

	// puede que el envido no haya sido cantado aun
	// lo deberia considerar aqui

	// tambien debo considerar el truco aqui
	if pdt.ElTrucoRespondible(p) {
		return b.jugar_truco(p), 1
	}

	// si no: elTrucoRespondible(p)
	return b.jugar_carta(p), 1
}
