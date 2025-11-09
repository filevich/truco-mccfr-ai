package info

import (
	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetRondaBaseFullBoolean struct {
	InfosetRondaBase
}

func (info *InfosetRondaBaseFullBoolean) setPuntos(p *pdt.Partida, m *pdt.Manojo) {
	our_team := m.Jugador.Equipo
	opp_team := m.Jugador.GetEquipoContrario()
	info.Nuestros_pts = 0
	info.Opp_pts = 0
	const threshold = 5

	if p.Puntajes[our_team] >= int(p.Puntuacion)-threshold {
		info.Nuestros_pts = 1
	}

	if p.Puntajes[opp_team] >= int(p.Puntuacion)-threshold {
		info.Opp_pts = 1
	}
}

func infosetRondaBaseFullBooleanFactory(
	a abs.IAbstraction,
) InfosetBuilder {
	return func(
		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,
	) Infoset {
		info := &InfosetRondaBaseFullBoolean{
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
