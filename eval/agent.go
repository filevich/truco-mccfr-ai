package eval

import (
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type Agent interface {
	Initialize()
	Free()
	UID() string
	Catch(*pdt.Partida, []enco.Envelope)
	ResetCatch()
	Action(p *pdt.Partida, inGameID string) (pdt.IJugada, float32)
}
