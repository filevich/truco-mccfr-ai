package cfr

import (
	"github.com/filevich/truco-mccfr-ai/utils"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type BotCFR struct {
	ID       string
	Filepath string
	Model    ITrainer
}

func (b *BotCFR) Initialize() {
	// lo cargo SOLO si no fue cargado aun
	if b.Model == nil {
		b.Model = LoadModel(b.Filepath, true, 1_000_000, false)
	}
}

func (b *BotCFR) Free() {
	b.Model = nil
}

func (b *BotCFR) SetUID(id string) {
	b.ID = id
}

func (b *BotCFR) UID() string {
	return b.ID
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
	// i := info.NewInfosetRondaBase(p, active_player, b.Model.GetAbs(), nil)
	i := b.Model.GetBuilder().Info(p, active_player, nil)
	hash, chi_len := i.Hash(b.Model.GetBuilder().Hash()), i.ChiLen()

	// obtengo la strategy
	strategy := b.Model.GetAvgStrategy(hash, chi_len)
	aix := utils.Sample(strategy)

	// obtengo el chi
	Chi := i.Iterable(p, active_player, aixs, b.Model.GetAbs())

	return Chi[aix], strategy[aix]
}

type BotCFR_Greedy struct {
	BotCFR
}

// Override only the Action method
func (b *BotCFR_Greedy) Action(
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
	i := b.Model.GetBuilder().Info(p, active_player, nil)
	hash, chi_len := i.Hash(b.Model.GetBuilder().Hash()), i.ChiLen()

	// obtengo la strategy
	strategy := b.Model.GetAvgStrategy(hash, chi_len)
	aix := utils.Argmax(strategy)

	// obtengo el chi
	Chi := i.Iterable(p, active_player, aixs, b.Model.GetAbs())

	return Chi[aix], strategy[aix]
}
