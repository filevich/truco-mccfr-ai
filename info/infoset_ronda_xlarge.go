package info

import (
	"github.com/filevich/truco-ai/abs"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type InfosetRondaXLarge struct {
	InfosetRondaLarge
}

func (info *InfosetRondaXLarge) setMuestra(p *pdt.Partida) {
	// InfosetRondaXLarge
	// info.muestra = 0
	// InfosetRondaXLarge
	info.muestra = p.Ronda.Muestra.Valor
	// InfosetRondaXXLarge
	// info.muestra = int(p.Ronda.Muestra.ID())
}

func infosetRondaXLargeFactory(

	a abs.IAbstraction,

) InfosetBuilder {

	return func(

		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,

	) Infoset {
		info := &InfosetRondaXLarge{}
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
