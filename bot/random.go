package bot

import (
	"math/rand"

	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type Random struct{}

func (b *Random) Initialize() {}

func (b *Random) Free() {}

func (b *Random) UID() string {
	return "Random"
}

func (b *Random) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *Random) ResetCatch() {}

// pre: el jugador no se fue al mazo
func (b *Random) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	m := p.Manojo(inGameID)
	chi := pdt.Chi(p, m)
	rix := rand.Intn(len(chi))

	// si no: elTrucoRespondible(p)
	return chi[rix], 1 / float32(len(chi))
}
