package cfr

import (
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type Agent interface {
	Initialize()
	Free()
	SetUID(string)
	UID() string
	Catch(*pdt.Partida, []enco.Envelope)
	ResetCatch()
	Action(p *pdt.Partida, inGameID string) (pdt.IJugada, float32)
}
