package info

import (
	"math"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetRondaBaseFullBucketed struct {
	InfosetRondaBase
}

func (info *InfosetRondaBaseFullBucketed) setPuntos(p *pdt.Partida, m *pdt.Manojo) {
	our_team := m.Jugador.Equipo
	opp_team := m.Jugador.GetEquipoContrario()
	const bucket = 5

	x := p.Puntajes[our_team]
	if x > int(p.Puntuacion) {
		x = int(p.Puntuacion)
	}
	nuestros_pts := int(math.Ceil(float64(x) / bucket))
	if nuestros_pts == 0 {
		nuestros_pts = 1
	}
	info.Nuestros_pts = nuestros_pts

	y := p.Puntajes[opp_team]
	if y > int(p.Puntuacion) {
		y = int(p.Puntuacion)
	}
	opp_pts := int(math.Ceil(float64(y) / bucket))
	if opp_pts == 0 {
		opp_pts = 1
	}
	info.Opp_pts = opp_pts
}

func infosetRondaBaseFullBucketedFactory(
	a abs.IAbstraction,
) InfosetBuilder {
	return func(
		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,
	) Infoset {
		info := &InfosetRondaBaseFullBucketed{
			InfosetRondaBase{
				Vision: m.Jugador.ID,
			},
		}
		chi_i := pdt.GetA(p, m)
		info.setMuestra(p)
		info.setNuestras_Cartas(p, m, a)
		info.setManojos_en_juego(p, m)
		info.setEnvido(p)
		info.setTruco(p)
		info.setChi(p, m, chi_i, a)
		info.setResultadoManos(p, m)
		info.setRonda(p, m, a)
		info.setPuntos(p, m)
		return info
	}
}
