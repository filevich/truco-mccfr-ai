package bot

import (
	"math"
	"math/rand"
	"sort"

	"github.com/filevich/truco-ai/info"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type dist struct {
	data  map[int]int
	total float32
}

func newDist(total float32, data map[int]int) *dist {
	return &dist{
		data:  data,
		total: total,
	}
}

// cummulative density function
func (d *dist) cdf(key int) float32 {
	s := 0
	for i := 0; i <= key; i++ {
		if v, ok := d.data[i]; ok {
			s += v
		}
	}
	return float32(s) / d.total
}

func (d *dist) probDareLineal(key int, lowerBoundDare float32) float32 {
	cdf := d.cdf(key)
	// y = mx + b
	// where `x` == `cdf`, `b` = `lowerBoundDare` and `m` = `1 - b`
	var (
		m = 1 - lowerBoundDare
		x = cdf
		b = lowerBoundDare
	)
	y := m*x + b
	return y
}

// alpha - hyperparam (e.g., alpha=3)
// k - correction constant (e.g., k=0.05)
func (d *dist) probDareTanh(key int, alpha, beta, gamma, k float32) float32 {
	// equation:
	// p = \frac{\tanh{(\alpha x - \frac{\alpha}{2} - \beta)}+1}{\gamma}+k
	var (
		x = d.cdf(key)
		z = alpha*x - alpha/2 - beta
		t = math.Tanh(float64(z))
	)
	y := (float32(t)+1)/gamma + k
	return y
}

type Pro struct {
	inGameID string
	// dists
	envidoDist *dist
	florDist   *dist
	powerDist  *dist
	// hypers
	Alphas []float32
	Betas  []float32
	Gammas []float32
	Ks     []float32
}

func (b *Pro) Initialize() {
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

func (b *Pro) Free() {}

func (b *Pro) UID() string {
	return "Pro"
}

func (b *Pro) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *Pro) ResetCatch() {
	// nueva ronda.
	// Aca es donde decido si voy a fingir demencia o no.
	// Hay dos juegos: el envite (envido+flor) y el truco.
	// Para ambos voy a decidir si juego o no en base a donde me encuentre en
	// las distribuciones de poder.

}

// retorna la carta mas chicas y mas poderosa en su poder SEGUN b.Abs.Abstraer(...) !
// si no tiene cartas en su poder, retorna -1,-1
func (b *Pro) cartas(p *pdt.Partida) (min, max *pdt.Carta) {
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

func (b *Pro) jugarCarta(p *pdt.Partida) pdt.IJugada {
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
		alpha, beta, gamma, k := b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k+0.1) < rand.Float32()
		_, ok = pdt.IrseAlMazo{JID: b.inGameID}.Ok(p)
		if stop || !ok {
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
func (b *Pro) getMaxFlor(p *pdt.Partida) int {
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
func (b *Pro) getMaxEnvite(p *pdt.Partida) int {
	maxEnvite := 0
	for i := 0; i < len(p.Ronda.Manojos); i++ {
		e := p.Ronda.Manojos[i].CalcularEnvido(p.Ronda.Muestra)
		if e > maxEnvite {
			maxEnvite = e
		}

	}
	return maxEnvite
}

func (b *Pro) getMax3ManojoSumPower(p *pdt.Partida, m *pdt.Manojo) int {
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

func (b *Pro) jugarFlor(p *pdt.Partida) pdt.IJugada {
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
		alpha   float32 = 0
		beta    float32 = 0
		gamma   float32 = 0
		k       float32 = 0
		maxFlor int     = b.getMaxFlor(p)
		stop    bool    = false
		ok      bool    = false
	)

	switch p.Ronda.Envite.Estado {
	case pdt.FLOR:
		// alpha, beta, gamma, k = 4, 0, 2, 0
		alpha, beta, gamma, k = b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.CantarFlor{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		// alpha, beta, gamma, k = 6, 0.9, 2, 0
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.CantarContraFlor{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.CantarFlor{
				JID: b.inGameID,
			}
		}
		// alpha, beta, gamma, k = 8, 2, 2, 0
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
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
		// alpha, beta, gamma, k = 6, 0.9, 2, 0
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
		if stop {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		// alpha, beta, gamma, k = 8, 2, 2, 0
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
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
		// alpha, beta, gamma, k = 8, 2, 2, 0
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.florDist.probDareTanh(maxFlor, alpha, beta, gamma, k) < rand.Float32()
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

func (b *Pro) jugarEnvido(p *pdt.Partida) pdt.IJugada {

	var (
		alpha     float32 = 0
		beta      float32 = 0
		gamma     float32 = 0
		k         float32 = 0
		maxEnvido int     = b.getMaxEnvite(p)
		stop      bool    = false
		ok        bool    = false
	)

	switch p.Ronda.Envite.Estado {
	case pdt.NOCANTADOAUN:
		alpha, beta, gamma, k = b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.TocarEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.TocarRealEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.TocarEnvido{
				JID: b.inGameID,
			}
		}
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
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
		alpha, beta, gamma, k = b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.TocarRealEnvido{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderQuiero{
				JID: b.inGameID,
			}
		}
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
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
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return &pdt.ResponderNoQuiero{
				JID: b.inGameID,
			}
		}
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
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
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.envidoDist.probDareTanh(maxEnvido, alpha, beta, gamma, k) < rand.Float32()
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
func (b *Pro) responderElTruco(p *pdt.Partida) pdt.IJugada {
	var (
		alpha  float32 = 0
		beta   float32 = 0
		gamma  float32 = 0
		k      float32 = 0
		cPower int     = b.getMax3ManojoSumPower(p, p.Manojo(b.inGameID))
		stop   bool    = false
		ok     bool    = false
	)

	switch p.Ronda.Truco.Estado {
	case pdt.TRUCO:
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.GritarReTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			// ok, ahora bajo el nivel y decido entre mazo o truco-querido
			alpha, beta, gamma, k = b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
			stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
			_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
			if stop || !ok {
				return &pdt.IrseAlMazo{
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
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.GritarVale4{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			// ok, ahora bajo el nivel y decido entre mazo o truco-querido
			alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
			stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
			_, ok = pdt.ResponderQuiero{JID: b.inGameID}.Ok(p)
			if stop || !ok {
				return &pdt.IrseAlMazo{
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
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
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
func (b *Pro) testearElTruco(p *pdt.Partida) pdt.IJugada {

	nuestro := p.Ronda.Truco.Estado > pdt.NOGRITADOAUN &&
		p.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == p.Manojo(b.inGameID).Jugador.Equipo
	if !nuestro {
		return nil
	}

	var (
		alpha  float32 = 0
		beta   float32 = 0
		gamma  float32 = 0
		k      float32 = 0
		cPower int     = b.getMax3ManojoSumPower(p, p.Manojo(b.inGameID))
		stop   bool    = false
		ok     bool    = false
	)

	switch p.Ronda.Truco.Estado {
	case pdt.NOGRITADOAUN:
		alpha, beta, gamma, k = b.Alphas[0], b.Betas[0], b.Gammas[0], b.Ks[0]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.GritarTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return nil
		}
		return &pdt.GritarTruco{
			JID: b.inGameID,
		}

	case pdt.TRUCOQUERIDO:
		alpha, beta, gamma, k = b.Alphas[1], b.Betas[1], b.Gammas[1], b.Ks[1]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
		_, ok = pdt.GritarReTruco{JID: b.inGameID}.Ok(p)
		if stop || !ok {
			return nil
		}
		return &pdt.GritarReTruco{
			JID: b.inGameID,
		}

	case pdt.RETRUCOQUERIDO:
		alpha, beta, gamma, k = b.Alphas[2], b.Betas[2], b.Gammas[2], b.Ks[2]
		stop = b.powerDist.probDareTanh(cPower, alpha, beta, gamma, k) < rand.Float32()
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
func (b *Pro) Action(

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
