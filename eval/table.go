package eval

import "github.com/filevich/truco-mccfr-ai/utils"

type Table map[string](map[string]*Results)

func NewTable() Table {
	return make(map[string](map[string]*Results))
}

func (tabla Table) Add(agent1, agent2 string, res *Results) {
	if _, ok := tabla[agent1]; !ok {
		tabla[agent1] = make(map[string]*Results)
	}
	tabla[agent1][agent2] = res
}

func (tabla Table) Metrics(agent1, agent2 string) (WP, ADP float32) {
	if _, ok := tabla[agent1][agent2]; ok {
		r := tabla[agent1][agent2]
		return r.WP(), r.ADP()
	}
	// debo devolver el inverso
	r := tabla[agent2][agent1]
	return 1 - r.WP(), -r.ADP()
}

func (tabla Table) WaldInterval(agent1, agent2 string) (upper, lower float64) {
	trials, success := 0, 0

	if _, ok := tabla[agent1][agent2]; ok {
		r := tabla[agent1][agent2]
		trials, success = r.TotalNumberOfGames, r.WonByACounter
	} else {
		// debo devolver el inverso
		r := tabla[agent2][agent1]
		trials, success = r.TotalNumberOfGames, r.TotalNumberOfGames-r.WonByACounter
	}

	return utils.WaldAdjusted(success, trials)
}
