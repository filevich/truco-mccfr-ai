package bot

import (
	"math/rand"

	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type BotRandom struct{}

func (b *BotRandom) Initialize() {}

func (b *BotRandom) Free() {}

func (b *BotRandom) UID() string {
	return "Random"
}

func (b *BotRandom) Catch(*pdt.Partida, []*enco.Envelope) {}

func (b *BotRandom) ResetCatch() {}

// pre: el jugador no se fue al mazo
func (b *BotRandom) Action(

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
