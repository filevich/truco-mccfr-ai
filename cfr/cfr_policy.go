package cfr

import (
	"github.com/truquito/bot/pers"
	"github.com/truquito/gotruco/pdt"
)

type CFR_Policy struct {
	Model *BotCFR
}

func (pol *CFR_Policy) Action(p *pers.Pers, mId string) pdt.IJugada {
	a, _ := pol.Model.Action(p.P, mId)
	return a
}

func (pol *CFR_Policy) Hash(p *pers.Pers, mId string) string {
	active_player := p.P.Manojo(mId)
	i := pol.Model.Model.GetBuilder().Info(p.P, active_player, nil)
	hash, _ := i.Hash(pol.Model.Model.GetBuilder().Hash), i.ChiLen()
	return hash
}

type CFR_Policy_Greedy struct {
	Model *BotCFR_Greedy
}

func (pol *CFR_Policy_Greedy) Action(p *pers.Pers, mId string) pdt.IJugada {
	a, _ := pol.Model.Action(p.P, mId)
	return a
}

func (pol *CFR_Policy_Greedy) Hash(p *pers.Pers, mId string) string {
	active_player := p.P.Manojo(mId)
	i := pol.Model.Model.GetBuilder().Info(p.P, active_player, nil)
	hash, _ := i.Hash(pol.Model.Model.GetBuilder().Hash), i.ChiLen()
	return hash
}
