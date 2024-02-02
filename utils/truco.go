package utils

import (
	"github.com/truquito/truco/enco"
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
