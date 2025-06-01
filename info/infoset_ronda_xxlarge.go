package info

import (
	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetRondaXXLarge struct {
	InfosetRondaLarge
}

func (info *InfosetRondaXXLarge) setMuestra(p *pdt.Partida) {
	// InfosetRondaLarge
	// info.muestra = 0
	// InfosetRondaXXLarge
	// info.muestra = p.Ronda.Muestra.Valor
	// InfosetRondaXXLarge
	info.muestra = int(p.Ronda.Muestra.ID())
}

func infosetRondaXXLargeFactory(

	a abs.IAbstraction,

) InfosetBuilder {

	return func(

		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,

	) Infoset {
		info := &InfosetRondaXXLarge{}
		info.setMuestra(p)
		info.setNumMano(p)
		info.setRixMe(p, m)
		info.setRixTurno(p)
		info.setManojosEnJuego(p)
		info.setNuestrasCartas(p, m, a)
		info.setTiradas(p, a)
		info.setHistory(p, msgs)
		info.setChi(p, m, a)
		return info
	}

}
