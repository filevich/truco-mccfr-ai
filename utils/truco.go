package utils

import (
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

func IsDoneAndPts(pkts []enco.Envelope) (bool, int, string) {
	done := false
	pts := -1
	autor := "-1"

	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TNuevaPartida ||
			pkt.Message.Cod() == enco.TNuevaRonda ||
			pkt.Message.Cod() == enco.TRondaGanada {
			done = true
		} else if pkt.Message.Cod() == enco.TSumaPts {
			m, _ := pkt.Message.(enco.SumaPts)
			pts = m.Puntos
			autor = m.Autor
		}
	}
	return done, pts, autor
}

func IsDoneAndPtsFull(pkts []enco.Envelope) (int, string) {
	pts := 0
	autor := "-1"

	for _, pkt := range pkts {
		if pkt.Message.Cod() == enco.TSumaPts {
			m, _ := pkt.Message.(enco.SumaPts)
			pts += m.Puntos
			autor = m.Autor
		}
	}
	return pts, autor
}

// indece del jugador relativo AL MANO
// mano --> turno -----> pie/respondedor-envite ---> yo
func RIX(p *pdt.Partida, m *pdt.Manojo) int {
	n := len(p.Ronda.Manojos)
	// return mod(p.Ronda.GetIdx(*m)-n, n)
	// return mod(m.Jugador.Jix-n, n)
	return Mod(p.Ronda.JIX(m.Jugador.ID)-n, n)
}

func Payoffs(elMano *pdt.Manojo, pts int, autor *pdt.Manojo) []float32 {
	if autor.Jugador.Equipo == elMano.Jugador.Equipo {
		return []float32{float32(pts), float32(-pts)}
	}

	return []float32{float32(-pts), float32(pts)}
}
