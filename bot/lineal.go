package bot

import (
	"math"
	"math/rand"
	"sort"

	"github.com/filevich/truco-ai/info"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

/*

0 flor/noquiero
1 contraflor/flor
2 contrafloralrest/contraflor
3 quiero_contraflor/noquiero_contraflor
4 contraflor_alrest/quiero
5 quiero_contrafloralresto/noquiero_contrafloralrest

6 envido/noquiero_envido
7 realenvido/envido
8 faltaenvido/realenvido

9 quiero_envido/noquiero_envido
10 real_envido/quiero_envido
11 falta_envid/real_envido

12 quiero_realenvido/noquiero_realenvido
13 faltaenvido/quiero_realevnido

14 quiero_faltaenvid/noquiero_faltaenvido

15 retruco/ (16 quiero_truco/noquiero_truco)
17 vale4/ (18 quiero_retruco/noquiero_retruco)
19 quiero_vale4/noquiero_vale4

20 truco/nada
21 retruco/nada
22 vale4/nada

23 mazo/seguir

*/

type Lineal struct {
	inGameID string
	// dists
	envidoDist *dist
	florDist   *dist
	powerDist  *dist
	// lower
	LowerBounds []float32
}

func (b *Lineal) Initialize() {
	if b.envidoDist == nil {
		b.envidoDist = newDist(
			308_800,
			map[int]int{0: 1512, 1: 2340, 10: 8752, 11: 7440,
				12: 7048, 13: 6348, 14: 3544, 2: 2572, 20: 8640, 21: 9000,
				22: 8640, 23: 11880, 24: 12480, 25: 14400, 26: 14760, 27: 21720,
				28: 15000, 29: 15240, 3: 4140, 30: 16080, 31: 14280, 32: 11520,
				33: 15840, 34: 13200, 35: 8280, 36: 6480, 37: 3360, 4: 4600,
				5: 5532, 6: 7516, 7: 9216, 8: 8872, 9: 8568,
			},
		)
	}

	if b.florDist == nil {
		b.florDist = newDist(
			56_760,
			map[int]int{0: 120, 1: 360, 10: 1236, 11: 1228, 12: 960,
				13: 848, 14: 508, 15: 360, 16: 268, 17: 120, 18: 120, 2: 360,
				27: 720, 28: 1104, 29: 1452, 3: 720, 30: 2076, 31: 2248, 32: 2564,
				33: 2828, 34: 3628, 35: 3668, 36: 3736, 37: 3928, 38: 3052,
				39: 2708, 4: 728, 40: 2316, 41: 1856, 42: 1572, 43: 1212, 44: 864,
				45: 392, 46: 236, 47: 40, 5: 1080, 6: 1200, 7: 1568, 8: 1328,
				9: 1448,
			},
		)
	}

	if b.powerDist == nil {
		b.powerDist = newDist(
			365_560,
			map[int]int{
				30: 40, 31: 360, 32: 828, 33: 1678, 34: 2214, 35: 3270, 36: 4426,
				37: 5886, 38: 7054, 39: 8808, 40: 10371, 41: 12052, 42: 12889,
				43: 13955, 44: 14705, 45: 15464, 46: 14937, 47: 14839, 48: 14184,
				49: 13444, 50: 12137, 51: 11317, 52: 10694, 53: 10027, 54: 9309,
				55: 8849, 56: 8837, 57: 8504, 58: 8389, 59: 8637, 60: 8913,
				61: 8827, 62: 8619, 63: 8512, 64: 8018, 65: 7361, 66: 6509,
				67: 5930, 68: 5046, 69: 4147, 70: 3268, 71: 2829, 72: 2229,
				73: 1911, 74: 1670, 75: 1639, 76: 1410, 77: 1356, 78: 1202,
				79: 1222, 80: 1142, 81: 1103, 82: 1042, 83: 922, 84: 744, 85: 549,
				86: 429, 87: 234, 88: 156, 89: 78, 90: 39, 93: 40, 94: 40, 95: 80,
				96: 80, 97: 80, 98: 40, 99: 40,
			},
		)
	}
}

func (b *Lineal) Free() {}

func (b *Lineal) UID() string {
	return "Lineal"
}

func (b *Lineal) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *Lineal) ResetCatch() {
	// nueva ronda.
	// Aca es donde decido si voy a fingir demencia o no.
	// Hay dos juegos: el envite (envido+flor) y el truco.
	// Para ambos voy a decidir si juego o no en base a donde me encuentre en
	// las distribuciones de poder.

}

// retorna la carta mas chicas y mas poderosa en su poder SEGUN b.Abs.Abstraer(...) !
// si no tiene cartas en su poder, retorna -1,-1
func (b *Lineal) cartas(p *pdt.Partida) (min, max *pdt.Carta) {
	minPoder, maxPoder := math.MaxInt32, math.MinInt32
	minCartaIx, maxCartaIx := -1, -1
	m := p.Manojo(b.inGameID)
	for cix, tirada := range m.Tiradas {
		if !tirada {
			poder := m.Cartas[cix].CalcPoder(p.Ronda.Muestra) // antes
			// poder := b.Abs.Abstraer(m.Cartas[cix], &p.Ronda.Muestra) // ahora
			if poder > maxPoder {
				maxPoder = poder
				maxCartaIx = cix
			}
			if poder < minPoder {
				minPoder = poder
				minCartaIx = cix
			}
		}
	}

	return m.Cartas[minCartaIx], m.Cartas[maxCartaIx]
}

func (b *Lineal) jugarCarta(p *pdt.Partida) pdt.IJugada {
	m := p.Manojo(b.inGameID)
	_, maxOpCarta, vamos := info.Vamos(p, m)

	switch vamos {
	case "ganando":
		// tiro la mas chica que tenga porque no hay necesidad
		minC, _ := b.cartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *minC,
		}

	// es indiferente la abstraccion porque:
	// el caso ganando/perdiendo/empatados se calcula segun Carta.CalcPoder(...)
	// la cual es una funcion independiente de la abstraccion

	case "perdiendo", "empatados", "?":

		// rinde irnos al mazo?
		cPower := b.getMax3ManojoSumPower(p, p.Manojo(b.inGameID))
		stop := false
		ok := false
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[23]) < rand.Float32()
		_, ok = pdt.IrseAlMazo{JID: b.inGameID}.Ok(p)
		if stop && ok {
			return &pdt.IrseAlMazo{
				JID: b.inGameID,
			}
		}

		// la idea es tirar la MINIMA carta TIRABLE en nuestro poder en esta mano
		// si la tengo yo -> entonces la tiro
		// si no la tengo yo -> entonces tiro la de menor poder

		// empezando por mi, me fijo de mi equipo quién tiene la carta más baja
		// que le gane a la `maxOpCarta`

		minPoderCartaSuperior := math.MaxInt
		esMia := false
		cartaIdx := -1

		for i := p.Ronda.MIXS[m.Jugador.ID]; i < len(p.Ronda.Manojos); i++ {
			manojo := p.Ronda.Manojos[i]
			if esDeMiEquipo := manojo.Jugador.Equipo == m.Jugador.Equipo; !esDeMiEquipo {
				continue
			}

			for cix, c := range manojo.Cartas {
				if manojo.Tiradas[cix] {
					continue
				}
				// no ha sido tirada.
				// es superior?
				poder := c.CalcPoder(p.Ronda.Muestra)
				maxOp := -1
				if maxOpCarta != nil {
					maxOp = maxOpCarta.CalcPoder(p.Ronda.Muestra)
				}
				if esSuperior := poder > maxOp; esSuperior {
					// ok, es menor de la mejor encontrada hasta ahora?
					if esMasEconomica := poder < minPoderCartaSuperior; esMasEconomica {
						minPoderCartaSuperior = poder
						cartaIdx = cix
						esMia = manojo.Jugador.ID == m.Jugador.ID
					}
				}
			}
		}

		if esMia {
			return &pdt.TirarCarta{
				JID:   m.Jugador.ID,
				Carta: *m.Cartas[cartaIdx],
			}
		}

		minC, _ := b.cartas(p)
		return &pdt.TirarCarta{
			JID:   m.Jugador.ID,
			Carta: *minC,
		}

	//
	default:
		return nil
	}
}

// retorna la maxima flor nuestra (que no se haya ido al mazo aun)
func (b *Lineal) getMaxFlor(p *pdt.Partida) int {
	maxFlor := 0
	for i := 0; i < len(p.Ronda.Manojos); i++ {
		if flor, err := p.Ronda.Manojos[i].CalcFlor(p.Ronda.Muestra); err != nil {
			if flor > maxFlor {
				maxFlor = flor
			}
		}
	}
	return maxFlor
}

// max envite
func (b *Lineal) getMaxEnvite(p *pdt.Partida) int {
	maxEnvite := 0
	for i := 0; i < len(p.Ronda.Manojos); i++ {
		e := p.Ronda.Manojos[i].CalcularEnvido(p.Ronda.Muestra)
		if e > maxEnvite {
			maxEnvite = e
		}

	}
	return maxEnvite
}

func (b *Lineal) getMax3ManojoSumPower(p *pdt.Partida, m *pdt.Manojo) int {
	numPlayers := len(p.Ronda.Manojos)
	cardPowers := make([]int, 0, 3*(numPlayers>>1))

	for _, manojo := range p.Ronda.Manojos {
		if isTeammate := manojo.Jugador.Equipo == m.Jugador.Equipo; !isTeammate {
			continue
		}
		if m.SeFueAlMazo {
			continue
		}
		for _, c := range manojo.Cartas {
			power := c.CalcPoder(p.Ronda.Muestra)
			cardPowers = append(cardPowers, power)
		}
	}

	sort.Ints(cardPowers)

	// sum up only top 3 cards
	sum := 0
	for _, v := range cardPowers[len(cardPowers)-3:] {
		sum += v
	}

	return sum
}

func (b *Lineal) jugarFlor(p *pdt.Partida) pdt.IJugada {
	m := p.Manojo(b.inGameID)

	// Si nadie ha cantado flor aún (estado < flor) o si el último que cantó fue amigo mío:
	// Yo canto flor
	nadieHaCantadoFlorAun := p.Ronda.Envite.Estado < pdt.FLOR
	elUltimoFueDeMiEquipo := !nadieHaCantadoFlorAun &&
		p.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo == m.Jugador.Equipo

	if nadieHaCantadoFlorAun || elUltimoFueDeMiEquipo {
		return &pdt.CantarFlor{
			JID: b.inGameID,
		}
	}

	// ok ya sé que tengo que “retrucar”
	// si:flor -> puedo responder: noquiero;flor;cf;cfr
	// si:cf -> puedo responder: noquiero;quiero;cfr
	// si:cfr -> puedo responder: noquiero;quiero

	// o sea, hay 3 niveles:
	// flor -> beta=0
	// contraflor -> beta=1
	// contrafloralresto -> beta=1.5

	var (
		maxFlor int  = b.getMaxFlor(p)
		stop    bool = false
		ok      bool = false
	)

	switch p.Ronda.Envite.Estado {
	case pdt.FLOR:
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[0]) < rand.Float32()
		_, ok = pdt.CantarFlor{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[1]) < rand.Float32()
		_, ok = pdt.CantarContraFlor{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.CantarFlor{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[2]) < rand.Float32()
		_, ok = pdt.CantarContraFlorAlResto{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.CantarContraFlor{
				JID: b.inGameID,
			}
		}
		return &pdt.CantarContraFlorAlResto{
			JID: b.inGameID,
		}

	case pdt.CONTRAFLOR:
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[3]) < rand.Float32()
		if stop {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[4]) < rand.Float32()
		_, ok = pdt.CantarContraFlorAlResto{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderQuiero{
				JID: b.inGameID,
			}
		}
		return &pdt.CantarContraFlorAlResto{
			JID: b.inGameID,
		}

	case pdt.CONTRAFLORALRESTO:
		stop = b.powerDist.probDareLineal(maxFlor, b.LowerBounds[5]) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		return &pdt.ResponderQuiero{
			JID: b.inGameID,
		}
	default:
		return nil
	}
}

func (b *Lineal) jugarEnvido(p *pdt.Partida) pdt.IJugada {

	var (
		maxEnvido int  = b.getMaxEnvite(p)
		stop      bool = false
		ok        bool = false
	)

	switch p.Ronda.Envite.Estado {
	case pdt.NOCANTADOAUN:
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[6]) < rand.Float32()
		_, ok = pdt.TocarEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[7]) < rand.Float32()
		_, ok = pdt.TocarRealEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.TocarEnvido{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[8]) < rand.Float32()
		_, ok = pdt.TocarFaltaEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.TocarRealEnvido{
				JID: b.inGameID,
			}
		}
		return &pdt.TocarFaltaEnvido{
			JID: b.inGameID,
		}

	case pdt.ENVIDO:
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[9]) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[10]) < rand.Float32()
		_, ok = pdt.TocarRealEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[11]) < rand.Float32()
		_, ok = pdt.TocarFaltaEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.TocarRealEnvido{
				JID: b.inGameID,
			}
		}
		return &pdt.TocarFaltaEnvido{
			JID: b.inGameID,
		}

	case pdt.REALENVIDO:
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[12]) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[13]) < rand.Float32()
		_, ok = pdt.TocarFaltaEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderQuiero{
				JID: b.inGameID,
			}
		}
		return &pdt.TocarFaltaEnvido{
			JID: b.inGameID,
		}

	case pdt.FALTAENVIDO:
		stop = b.powerDist.probDareLineal(maxEnvido, b.LowerBounds[14]) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		return &pdt.ResponderQuiero{
			JID: b.inGameID,
		}
	default:
		return nil
	}
}

// siempre responde; nunca retorna `nil`
func (b *Lineal) responderElTruco(p *pdt.Partida) pdt.IJugada {
	var (
		cPower int  = b.getMax3ManojoSumPower(p, p.Manojo(b.inGameID))
		stop   bool = false
		ok     bool = false
	)

	switch p.Ronda.Truco.Estado {
	case pdt.TRUCO:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[15]) < rand.Float32()
		_, ok = pdt.GritarReTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			// ok, ahora bajo el nivel y decido entre mazo o truco-querido
			stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[16]) < rand.Float32()
			_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
			if stop || !ok {
				return &pdt.ResponderNoQuiero{
					JID: b.inGameID,
				}
			} else {
				return &pdt.ResponderQuiero{
					JID: b.inGameID,
				}
			}
		}
		return &pdt.GritarReTruco{
			JID: b.inGameID,
		}

	case pdt.RETRUCO:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[17]) < rand.Float32()
		_, ok = pdt.GritarVale4{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			// ok, ahora bajo el nivel y decido entre mazo o truco-querido
			stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[18]) < rand.Float32()
			_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
			if stop || !ok {
				return &pdt.ResponderNoQuiero{
					JID: b.inGameID,
				}
			} else {
				return &pdt.ResponderQuiero{
					JID: b.inGameID,
				}
			}
		}
		return &pdt.GritarVale4{
			JID: b.inGameID,
		}

	case pdt.VALE4:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[19]) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		return &pdt.ResponderQuiero{
			JID: b.inGameID,
		}
	default:
		return nil
	}
}

// está en no-gritado o el querido lo tengo yo
func (b *Lineal) testearElTruco(p *pdt.Partida) pdt.IJugada {

	nuestro := p.Ronda.Truco.Estado > pdt.NOGRITADOAUN &&
		p.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == p.Manojo(b.inGameID).Jugador.Equipo
	if !nuestro {
		return nil
	}

	var (
		cPower int  = b.getMax3ManojoSumPower(p, p.Manojo(b.inGameID))
		stop   bool = false
		ok     bool = false
	)

	switch p.Ronda.Truco.Estado {
	case pdt.NOGRITADOAUN:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[20]) < rand.Float32()
		_, ok = pdt.GritarTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return nil
		}
		return &pdt.GritarTruco{
			JID: b.inGameID,
		}

	case pdt.TRUCOQUERIDO:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[21]) < rand.Float32()
		_, ok = pdt.GritarReTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return nil
		}
		return &pdt.GritarReTruco{
			JID: b.inGameID,
		}

	case pdt.RETRUCOQUERIDO:
		stop = b.powerDist.probDareLineal(cPower, b.LowerBounds[22]) < rand.Float32()
		_, ok = pdt.GritarVale4{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return nil
		}
		return &pdt.GritarVale4{
			JID: b.inGameID,
		}
	// vale 4 querido ya no tiene "upgrade"
	default:
		return nil
	}
}

// pre: el jugador no se fue al mazo
func (b *Lineal) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	b.inGameID = inGameID

	if pdt.LaFlorEsRespondible(p) {
		if _, err := p.Manojo(b.inGameID).CalcFlor(p.Ronda.Muestra); err == nil {
			return b.jugarFlor(p), 1
		}
	}

	if pdt.ElEnvidoEsRespondible(p) {
		return b.jugarEnvido(p), 1
	}

	// tal vez el envido no es respondible pero si lo puedo iniciar
	if p.Ronda.Envite.Estado != pdt.DESHABILITADO {
		if j := b.jugarEnvido(p); j.ID() != pdt.JID_NO_QUIERO {
			return j, 1
		}
	}

	// tambien debo considerar el truco aqui
	if pdt.ElTrucoRespondible(p) {
		return b.responderElTruco(p), 1
	}

	if j := b.testearElTruco(p); j != nil {
		return j, 1
	}

	// debo considerar irme al mazo
	return b.jugarCarta(p), 1
}
