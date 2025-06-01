package dumbo

import (
	"fmt"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/pdt"
)

func esElUltimoEnTirar(p *pdt.Partida, m *pdt.Manojo) bool {
	return p.Ronda.GetSigHabilitado(*m) == nil
}

func maxTiradas(p *pdt.Partida) map[pdt.Equipo]*pdt.CartaTirada {

	var (
		maxPoder = map[pdt.Equipo]int{pdt.Rojo: -1, pdt.Azul: -1}
		max      = map[pdt.Equipo]*pdt.CartaTirada{pdt.Rojo: nil, pdt.Azul: nil}
		tiradas  = p.Ronda.GetManoActual().CartasTiradas
	)

	for i, tirada := range tiradas {
		poder := tirada.Carta.CalcPoder(p.Ronda.Muestra)
		equipo := p.Ronda.Manojo(tirada.Jugador).Jugador.Equipo
		if poder > maxPoder[equipo] {
			maxPoder[equipo] = poder
			max[equipo] = &tiradas[i]
		}
	}

	return max
}

// PRE: m es el utlimo habilitado en tirar
func hayMenorYSigueGanando(

	p *pdt.Partida,
	m *pdt.Manojo,
	tc pdt.TirarCarta,
	a abs.IAbstraction,

) bool {

	// fix null abs
	if a.String() == string(abs.NULL_ID) {
		a = &abs.A3{}
	}

	var (
		maxTiradas                 = maxTiradas(p)
		mejor_ellos                = maxTiradas[m.Jugador.GetEquipoContrario()]
		bucket_la_que_pienso_tirar = a.Abstract(&tc.Carta, &p.Ronda.Muestra)
		bucket_la_mejor_de_ellos   = a.Abstract(&mejor_ellos.Carta, &p.Ronda.Muestra)
	)

	// tengo alguna de MENOR bucket que aun asi le gana a la mejor tirada?
	hay_menor_y_gana := false
	for cix, tirada := range m.Tiradas {
		if !tirada {
			c := m.Cartas[cix]
			bucket := a.Abstract(c, &p.Ronda.Muestra)
			es_menor := bucket < bucket_la_que_pienso_tirar
			le_gana := bucket_la_mejor_de_ellos < bucket
			if es_menor && le_gana {
				hay_menor_y_gana = true
				break
			}
		}
	}

	return hay_menor_y_gana

	// d1: suponiendo que NO es parte de una estrategia (eg, se quiere "ganar el
	// turno" en la siguiente mano y por eso tira una carta "mas alta" o bien
	// quiere dar "una falsa idea de poderio".)

	// si ibamos GANANDO -> no habia necesidad de tirar una carta mas alta
	// (porque el win ya lo tenemos asegurado) no, porque podria ser una
	// estrategia
	// si vamos PERDIENDO -> no habia necesidad de tirar una carta mas alta
}

func _maxTiradasStr(p *pdt.Partida, abs abs.IAbstraction) string {
	tiradas := maxTiradas(p)
	s := "\n"
	for e, t := range tiradas {
		s += fmt.Sprintf("[%s] - %s - %s - %s ~> bucket #%d\n",
			e,
			t.Jugador,
			t.Carta,
			abs,
			abs.Abstract(&t.Carta, &p.Ronda.Muestra))
	}

	if tiradas[pdt.Azul].CalcPoder(p.Ronda.Muestra) > tiradas[pdt.Rojo].CalcPoder(p.Ronda.Muestra) {
		s += "va ganando azul\n"
	} else {
		s += "va ganando rojo\n"
	}

	return s
}

func IsDumbo(

	p *pdt.Partida,
	m *pdt.Manojo,
	j pdt.IJugada,
	a abs.IAbstraction,

) bool {

	// retorna true sii
	// 1. tira carta
	// 2. es el ultimo en tirar
	// 3. tira "sobrado"

	if j.ID() == pdt.JID_TIRAR_CARTA {
		if tc, ok := j.(*pdt.TirarCarta); ok {
			return esElUltimoEnTirar(p, m) && hayMenorYSigueGanando(p, m, *tc, a)
		}
		tc := j.(pdt.TirarCarta)
		return esElUltimoEnTirar(p, m) && hayMenorYSigueGanando(p, m, tc, a)

	}

	return false
}
