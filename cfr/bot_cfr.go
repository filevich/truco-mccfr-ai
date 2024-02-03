package cfr

import (
	"crypto/sha1"
	"strings"

	"github.com/filevich/truco-cfr/info"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type BotCFR struct {
	N     string
	F     string
	Model ITrainer
}

func (b *BotCFR) Initialize() {
	// lo cargo SOLO si no fue cargado aun
	if b.Model == nil {
		if strings.HasSuffix(b.F, ".json") {
			b.Model = Load(CFR_T, b.F)
		} else {
			b.Model = Load_model(b.F, true, 1_000_000)
		}
	}
}

func (b *BotCFR) Free() {
	b.Model = nil
}

func (b *BotCFR) UID() string {
	return b.N
}

func (b *BotCFR) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *BotCFR) ResetCatch() {}

// pre: el jugador no se fue al mazo
func (b *BotCFR) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	// pseudo jugador activo
	active_player := p.Manojo(inGameID)

	// obtengo el infoset
	aixs := pdt.GetA(p, active_player)
	i := info.NewInfosetRondaBase(p, active_player, b.Model.Get_abs(), nil)
	hash, chi_len := i.Hash(sha1.New()), i.ChiLen()

	// obtengo la strategy
	strategy := b.Model.Get_avg_strategy(hash, chi_len)
	aix := utils.Sample(strategy)

	// obtengo el chi
	Chi := i.Iterable(p, active_player, aixs, b.Model.Get_abs())

	return Chi[aix], strategy[aix]
}
