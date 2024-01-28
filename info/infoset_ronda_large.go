package info

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"sort"
	"strconv"

	"github.com/filevich/truco-cfr/abs"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

// S, SE, C, R, plus pro max ultra

// `InfosetRondaLarge`, Guarda toda la info de la ronda.
//
// Pero no guarda:
//   - ninguna informacion relacionada a la partida. I.e., puntaje y puntuacion
//   - "el historial" de jugadas. e.g., Bob:Mazo Ben:1C ~ Ben:1C Bob:Mazo
//   - ^ Debería considerarse como infosets diferentes?
//
// Notar que a diferencia de Infoset1 NO se guarda explicitamente si vamos
// "ganando"/"perdiendo"/"empatados"/"?". Esta info "es deducible" o "está
// embebida" en el campo `historial`

type InfosetRondaLarge struct {
	// 1. `muestra`
	// la almacena de forma pura; sin abstracción.
	// no es posible abstraer la muestra porque la función de abstraer depende
	// de la muestra misma.
	// notar que almacenar solor el valor de la muestra (i.e., el número)
	// también es un tipo de abstracción.
	muestra int

	// 2. `num_mano_actual`: int
	numMano int

	// 3. `rixMe` ~ RIX: who?
	// Mi posicion o indice relativo a la posicion de `elMano` actual
	// Obs: si `rixMe == 0` entonces yo soy `elMano`, si `rixMe == n - 1`
	// entonces yo soy el último jugador habilitado.
	// Notar que si `rixMe` es padr -> mi equipo es el poseedor de `elMano`.
	rixMe int

	// 4. `turno` ~ RIX who?
	// posicion relativa de `elTurno` respecto a `elMano`
	rixTurno int

	// 5. `ManojosEnJuego` quiénes se fueron al mazo y quiénes siguen en pie?
	// el indice 0 se corresponde con `elMano`
	manojosEnJuego []bool

	// 6. `nuestrasCartas` representa nuestras cartas.
	// Por cada jugador/manojo de nuestro equipo se almacena un único entero.
	// Estos enteros se calculan como la multiplicación de todos los primos
	// asociados a cada carta del manojo.
	// De esta forma, cartas diferentes generan identificadores diferentes
	// pero el orden (permutación) no altera el resultado.
	// El indice 0 se corresponde con las cartas del manojo más cercano de nuestro
	// equipo al Mano de la ronda actual.
	nuestrasCartas []int
	_miManojoPID   int // <- ignorar. cache solo para facilitar el código

	// 7. tiradas: [(abs(carta),who?)]
	tiradasCartas [][]int
	tiradasWho    [][]int

	// 8. historial (envido+flor+truco)
	// los cantos sobre cuánto tiene cada uno deben ser en orden
	historialQue    []string
	historialQuien  []int
	historialCuanto []int

	// 9. las acciones que YO puedo tomar
	// Las acciones de tipo `TirarCarta` son retornadas en orden según el valor
	// retornado por la abstracción de las cartas.
	// Notar también que si el jugador tiene en su posesión dos o más cartas
	// que comparten una misma abstracción/bucket, entonces a los efectos de la
	// abstracción, estas cartas son indistinguibles. Es por esto que en estos
	// casos se retorna (como máximo) solo una `IJugada` por cada bucket
	// posible.
	chi []pdt.IJugada
}

func (info *InfosetRondaLarge) setMuestra(p *pdt.Partida) {
	info.muestra = int(p.Ronda.Muestra.ID())
}

func (info *InfosetRondaLarge) setNumMano(p *pdt.Partida) {
	info.numMano = int(p.Ronda.ManoEnJuego)
}

func (info *InfosetRondaLarge) setRixMe(p *pdt.Partida, m *pdt.Manojo) {
	info.rixMe = RIX(p, m)
}

func (info *InfosetRondaLarge) setRixTurno(p *pdt.Partida) {
	info.rixTurno = RIX(p, p.Ronda.GetElTurno())
}

func (info *InfosetRondaLarge) setManojosEnJuego(p *pdt.Partida) {
	// tengo que empezar a iterar a partir del JIX del MANO
	n := len(p.Ronda.Manojos)
	info.manojosEnJuego = make([]bool, n)
	m := p.Ronda.GetElMano()
	for i := 0; i < n; i++ {
		info.manojosEnJuego[i] = !m.SeFueAlMazo
		m = p.Ronda.GetSiguiente(*m)
	}
}

func (info *InfosetRondaLarge) setNuestrasCartas(

	p *pdt.Partida,
	m *pdt.Manojo,
	a abs.IAbstraccion,

) {
	// cada equipo tiene n/2 manojos donde n es la cantidad de jugadores en la
	// partida.
	// Nota: dividir entre 2 es igual a hacer un shift right `n >> 1`
	n := len(p.Ronda.Manojos)
	info.nuestrasCartas = make([]int, 0, n>>1)
	e := m.Jugador.Equipo

	// tengo que empezar a iterar a partir del JIX del MANO
	manojo := p.Ronda.GetElMano()
	for i := 0; i < n; i++ {
		if esDeNuestroEquipo := manojo.Jugador.Equipo == e; esDeNuestroEquipo {
			pid := PrimifyManojo(manojo, &p.Ronda.Muestra, a)
			info.nuestrasCartas = append(info.nuestrasCartas, pid)
			if itsMe := manojo.Jugador.ID == m.Jugador.ID; itsMe {
				if !m.SeFueAlMazo {
					info._miManojoPID = pid
				} else {
					info._miManojoPID = 1
				}
			}
		}
		manojo = p.Ronda.GetSiguiente(*manojo)
	}
}

func (info *InfosetRondaLarge) setTiradas(p *pdt.Partida, a abs.IAbstraccion) {
	cartasTiradasPorMano := make([][]int, 3)
	whoTiradasPorMano := make([][]int, 3)

	for mix, mano := range p.Ronda.Manos {
		cantTiradas := len(mano.CartasTiradas)
		cartasTiradasPorMano[mix] = make([]int, cantTiradas)
		whoTiradasPorMano[mix] = make([]int, cantTiradas)
		for tix, tirada := range mano.CartasTiradas {
			cartasTiradasPorMano[mix][tix] = a.Abstraer(&tirada.Carta, &p.Ronda.Muestra)
			whoTiradasPorMano[mix][tix] = RIX(p, p.Manojo(tirada.Jugador))
		}
	}

	info.tiradasCartas = cartasTiradasPorMano
	info.tiradasWho = whoTiradasPorMano
}

func _esMsgHistory(msg enco.IMessage) bool {
	c := msg.Cod()

	if c == enco.TDiceSonBuenas ||
		c == enco.TCantarFlor ||
		c == enco.TCantarContraFlor ||
		c == enco.TCantarContraFlorAlResto ||
		c == enco.TTocarEnvido ||
		c == enco.TTocarRealEnvido ||
		c == enco.TTocarFaltaEnvido ||
		c == enco.TGritarTruco ||
		c == enco.TGritarReTruco ||
		c == enco.TGritarVale4 ||
		c == enco.TNoQuiero ||
		c == enco.TConFlorMeAchico ||
		c == enco.TQuieroTruco ||
		c == enco.TQuieroEnvite ||
		c == enco.TDiceTengo ||
		c == enco.TDiceSonMejores {
		return true
	}

	return false
}

func _parse(p *pdt.Partida, msg enco.IMessage) (

	quien int,
	que string,
	cuanto int,

) {

	que = string(msg.Cod())
	cuanto = 0

	switch enco.CodMsg(que) {

	case enco.TDiceSonBuenas:
		m, _ := msg.(enco.DiceSonBuenas)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TCantarFlor:
		m, _ := msg.(enco.CantarFlor)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TCantarContraFlor:
		m, _ := msg.(enco.CantarContraFlor)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TCantarContraFlorAlResto:
		m, _ := msg.(enco.CantarContraFlorAlResto)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TTocarEnvido:
		m, _ := msg.(enco.TocarEnvido)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TTocarRealEnvido:
		m, _ := msg.(enco.TocarRealEnvido)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TTocarFaltaEnvido:
		m, _ := msg.(enco.TocarFaltaEnvido)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TGritarTruco:
		m, _ := msg.(enco.GritarTruco)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TGritarReTruco:
		m, _ := msg.(enco.GritarReTruco)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TGritarVale4:
		m, _ := msg.(enco.GritarVale4)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TNoQuiero:
		m, _ := msg.(enco.NoQuiero)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TConFlorMeAchico:
		m, _ := msg.(enco.ConFlorMeAchico)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TQuieroTruco:
		m, _ := msg.(enco.QuieroTruco)
		quien = RIX(p, p.Manojo(string(m)))

	case enco.TQuieroEnvite:
		m, _ := msg.(enco.QuieroEnvite)
		quien = RIX(p, p.Manojo(string(m)))

	// (string, int)
	case enco.TDiceTengo:
		m, _ := msg.(enco.DiceTengo)
		quien = RIX(p, p.Manojo(m.Autor))
		cuanto = m.Valor

	case enco.TDiceSonMejores:
		m, _ := msg.(enco.DiceSonMejores)
		quien = RIX(p, p.Manojo(m.Autor))
		cuanto = m.Valor
	}

	return quien, que, cuanto
}

func (info *InfosetRondaLarge) setHistory(p *pdt.Partida, msgs []enco.IMessage) {
	info.historialQuien = make([]int, 0, len(msgs))
	info.historialQue = make([]string, 0, len(msgs))
	info.historialCuanto = make([]int, 0, len(msgs))

	for _, msg := range msgs {
		if _esMsgHistory(msg) {
			quien, que, cuanto := _parse(p, msg)
			info.historialQuien = append(info.historialQuien, quien)
			info.historialQue = append(info.historialQue, que)
			info.historialCuanto = append(info.historialCuanto, cuanto)
		}
	}
}

func (info *InfosetRondaLarge) setChi(

	p *pdt.Partida,
	m *pdt.Manojo,
	a abs.IAbstraccion,

) {
	// Notar que el simulador puede repartirme las mismas cartas pero con una
	// permutación diferente. E.g., [3c,1b,7c] v. [1b,3c,7c]
	// En ese caso, (y debido a que los manojos son representados como la mult.
	// de los primos asociados a las abstracciones de las cartas) ambos infosets
	// van a ser indistinguibles.
	// Ahora, si yo naivamente retorno `pdt.Chi` entoces a veces voy a estar
	// retornando [3c,1b,7c] y otras veces [1b,3c,7c], por lo que los regrets
	// y strategies se van a estar almacenando de forma aleatoria.
	// Es necesario, por lo tanto, retornar o bien las acciones ordenadas y
	// agrupadas según su abstracción

	// Supuestos importantes:
	//  1. la función `pdt.Chi` retorna primero las acciones de tipo
	//     `TirarCarta` y luego el resto
	//  2. el simulador NO ordena las cartas dentro del manojo y las "carga"
	//    de forma aleatoria

	// Para las cartas, voy a retornar un slice de `TirarCarta`, de tamaño
	// máximo 3, ordenas desde el bucket/abstracción con menor numeración hasta
	// la de mayor numeración.
	chi := pdt.Chi(p, m)

	bucketsSeen := make(map[int]pdt.IJugada)
	buckets := make([]int, 0, 3) // como máximo son 3 cartas; indep. de la abs.

	n := len(chi)
	ixAccionesStart := n

	for i, jugada := range chi {
		if jugada.ID() == pdt.JID_TIRAR_CARTA {
			tirar, _ := jugada.(pdt.TirarCarta)
			bucket := a.Abstraer(&tirar.Carta, &p.Ronda.Muestra)
			if _, ok := bucketsSeen[bucket]; ok {
				continue
			}
			bucketsSeen[bucket] = tirar
			buckets = append(buckets, bucket)
		} else {
			ixAccionesStart = i
			break // supuesto 1.
			// ya no hay más jugadas de tipo `TirarCarta`
		}
	}

	cantBucketsParaTirar := len(bucketsSeen)
	if cantBucketsParaTirar == 0 {
		info.chi = chi
		return
	}

	// ahora las agrego en orden según su bucket
	res := make([]pdt.IJugada, 0, n)
	sort.Ints(buckets)
	for _, bucket := range buckets {
		res = append(res, bucketsSeen[bucket])
	}

	info.chi = append(res, chi[ixAccionesStart:]...)
}

func (info *InfosetRondaLarge) ChiLen() int {
	return len(info.chi)
}

func (info *InfosetRondaLarge) Dump(_ bool) string {
	return fmt.Sprintf("%+v", info)
}

func (info *InfosetRondaLarge) HashBytes(h hash.Hash) []byte {
	hsep := []byte(sep)

	// 1. muestra int
	h.Write([]byte(strconv.Itoa(info.muestra)))
	h.Write(hsep)

	// 2. numMano int
	h.Write([]byte(strconv.Itoa(info.numMano)))
	h.Write(hsep)

	// 3. rixMe int
	h.Write([]byte(strconv.Itoa(info.rixMe)))
	h.Write(hsep)

	// 4. rixTurno int
	h.Write([]byte(strconv.Itoa(info.rixMe)))
	h.Write(hsep)

	// 5. manojosEnJuego []bool
	{
		bs, _ := json.Marshal(info.manojosEnJuego)
		h.Write([]byte(bs))
		h.Write(hsep)
	}

	// 6. nuestrasCartas []int
	{
		bs, _ := json.Marshal(info.nuestrasCartas)
		h.Write([]byte(bs))
		h.Write(hsep)
	}

	// 7.
	// tiradasCartas [][]int
	// tiradasWho    [][]int
	{
		bs, _ := json.Marshal(info.tiradasCartas)
		h.Write([]byte(bs))
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.tiradasWho)
		h.Write([]byte(bs))
		h.Write(hsep)
	}

	// 8.
	// historialQuien  []int
	// historialQue    []string
	// historialCuanto []int
	{
		bs, _ := json.Marshal(info.historialQuien)
		h.Write([]byte(bs))
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.historialQue)
		h.Write([]byte(bs))
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.historialCuanto)
		h.Write([]byte(bs))
		h.Write(hsep)
	}

	// 9. chi []pdt.IJugada
	{
		chi := make([]int, 0, len(info.chi))
		// El primer indice lo uso para almacenar un "_miManojoPID"
		// Es decir, un número el cual queda definido como la multiplicación de
		// todos los primos correspondientes luego de aplicar la abstracción
		// sobre las cartas que tengo disponibles para tirar.
		// En caso de que no tenga niguna carta para tirar, este indice queda
		// con el neutro de la multiplicación (i.e., el 1).
		chi = append(chi, info._miManojoPID)

		// Ahora agrego las otras jugadas/acciones disponibles a la derecha de
		// este número
		for _, j := range info.chi {
			if j.ID() != pdt.JID_TIRAR_CARTA {
				chi = append(chi, int(j.ID()))
			}
		}

		bs, _ := json.Marshal(chi)
		h.Write([]byte(bs))
		// h.Write(hsep) // <- not necessary
	}

	return h.Sum(nil)
}

func (info *InfosetRondaLarge) Hash(h hash.Hash) string {
	return hex.EncodeToString(info.HashBytes(h))
}

func (info *InfosetRondaLarge) Iterable(

	p *pdt.Partida,
	m *pdt.Manojo,
	_ pdt.A,
	a abs.IAbstraccion,

) []pdt.IJugada {

	return info.chi
}

func NewInfosetRondaLarge(

	p *pdt.Partida,
	m *pdt.Manojo,
	a abs.IAbstraccion,
	msgs []enco.IMessage,

) Infoset {

	info := &InfosetRondaLarge{}

	info.setMuestra(p)
	info.setNumMano(p)
	info.setRixMe(p, m)
	info.setRixTurno(p)
	info.setManojosEnJuego(p)
	info.setNuestrasCartas(p, m, a)
	info.setTiradas(p, a)
	info.setHistory(p, msgs)
	info.setChi(p, m, a)

	return info
}
