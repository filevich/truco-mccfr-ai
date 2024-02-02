package eval

import "github.com/filevich/truco-cfr/utils"

type Tabla map[string](map[string]*Resultados)

func New_Tabla() Tabla {
	return make(map[string](map[string]*Resultados))
}

func (tabla Tabla) Add(agent1, agent2 string, res *Resultados) {
	if _, ok := tabla[agent1]; !ok {
		tabla[agent1] = make(map[string]*Resultados)
	}
	tabla[agent1][agent2] = res
}

func (tabla Tabla) Metrics(agent1, agent2 string) (WP, ADP float32) {
	if _, ok := tabla[agent1][agent2]; ok {
		r := tabla[agent1][agent2]
		return r.WP(), r.ADP()
	}
	// debo devolver el inverso
	r := tabla[agent2][agent1]
	return 1 - r.WP(), -r.ADP()
}

func (tabla Tabla) WaldInterval(agent1, agent2 string) (upper, lower float64) {
	trials, success := 0, 0

	if _, ok := tabla[agent1][agent2]; ok {
		r := tabla[agent1][agent2]
		trials, success = r.Total, r.Wins_by_A_count
	} else {
		// debo devolver el inverso
		r := tabla[agent2][agent1]
		trials, success = r.Total, r.Total-r.Wins_by_A_count
	}

	return utils.WaldAdjusted(success, trials)
}

// func (tabla Tabla) NormalInterval(agent1, agent2 string) (upper, lower float64) {
// 	var data []float64 = nil

// 	if _, ok := tabla[agent1][agent2]; ok {
// 		data = tabla[agent1][agent2].Wins
// 	} else {
// 		n := len(tabla[agent2][agent1].Wins)
// 		inverse := make([]float64, n)
// 		for i, v := range tabla[agent2][agent1].Wins {
// 			inverse[i] = 1 - v
// 		}
// 		data = inverse
// 	}

// 	loc, _ := stats.Median(data)
// 	scale, _ := stats.StdDevS(data)
// 	alpha := 0.1
// 	i := stats.NormInterval(alpha, loc, scale)

// 	return i[1], i[0]
// }
