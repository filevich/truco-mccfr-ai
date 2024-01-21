package info

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"math/rand"
	"strconv"

	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

/*

Primera de implementación de `Infoset` propuesta en 2022 durante el PDG.
Solo abstrae la `Ronda` mas no la `Partida`.

Lo bueno de esta implementación:
 - Es "liviana"/"barata" y a los efectos prácticos obtiene buenos resultados
 - Usa sha1 -> 160 bit max -> 2^160 = 1.46e48 max

Lo malo de esta implementación:
 - Solo se abstrae información relacionada a la Ronda; nada de la Partida
 - No es compatible con la abstracción `Null`
 - Descarta la info sobre qué cartas tiraron OP en las manos anteriores
 - No contiene la info sobre cómo va la puntuación a nivel de parida; esto hace
   que el "falta-envido" o "contra-flor-al-rest" no se puedan "calcular".

*/

type InfosetRondaBase struct {
	// debug
	Vision string

	// De la muestra solo se almacena el valor, NO el palo. Esto hace que la
	// abstraccion `Null` no sea compatible con esta implementacion ya que un
	// jugador no podria distinguir si tiene o no una pieza en su manojo.
	Muestra int

	// `NuestrasCartas` contiene las abstracciones de todas nuestras cartas
	// solo las cartas NO TIRADAS
	NuestrasCartas [][]int

	// `ManojosEnJuego` almacena `true` por cada jugador que aun no se fue al
	// mazo
	ManojosEnJuego []bool

	// Notar que esta implementacion de `Infoset` se pierde la informacion sobre
	// qué cartas tiraron los oponentes en las manos anteriores.
	// La informacion a tener en cuenta sobre las manos de la ronda actual se
	// limita a "el resultado de cada mano" (`Resultado_manos`) y cómo vamos en la
	// mano actual (`Mano_actual`).

	// puede ser "en-juego", "parda", "ganada", "perdida"
	ResultadoManos []string

	// cómo vamos en esta mano?
	ManoActual struct {
		Max_us int
		Max_op int
		Vamos  string
	}

	// estado de los "subjuegos"
	// notar que si el envido se encuentra en "ENVIDO" y yo tengo las acciones:
	// ["envido", "real-envido", "falta-envido", "quiero", "no-quiero"]
	// entonces se puede deducir que el que lo tocó fue el equipo contario.
	// lo mismo para el truco
	Envido string
	Truco  string

	// las acciones que YO puedo tomar
	Chi []int
}

func (info *InfosetRondaBase) setMuestra(p *pdt.Partida) {
	// 1. muestra:
	//	el numero/valor de la muestra (no se inculye el palo)

	//	formato:
	//		valor--
	info.Muestra = p.Ronda.Muestra.Valor
}

func (info *InfosetRondaBase) setNuestras_Cartas(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	abs abs.IAbstraccion,

) {
	// 2. cartas "nuestras":

	// por ejemplo, si usaramos la abstraccion A1, las cartas
	// serian agrupadas en 3 conjuntos:
	//  {2,4,5,11,10}, {1,1,7,7}, {3,2,1,12,11,10,7,6,5,4}

	// si jugamos de a 4p, esta funcion genera un slice 2d del tipo:
	// [[2,1,1],[3,0]] // <- solo las cartas no tiradas

	// luego, la funcion de hash se encarga de fusionar estos slices:
	// como solo importa la combinacion y no la permutacion
	// asocio un primo a cada agrupacion y multiplico estos primos

	// formato:
	// p1.p2.p3--
	//
	// 		donde:
	// 			p1 = prime(c11) * prime(c12) * prime(c13)
	//    	p2 = prime(c21) * prime(c22) * prime(c23)
	// 			p3 = prime(c31) * prime(c31) * prime(c32)

	info.NuestrasCartas = make([][]int, len(p.Ronda.Manojos)/2)
	m := manojo
	for i := 0; i < len(p.Ronda.Manojos); i++ {
		mismoEquipo := m.Jugador.Equipo == manojo.Jugador.Equipo
		if mismoEquipo {
			info.NuestrasCartas[i/2] = make([]int, 0, 3)
			for cix, c := range m.Cartas {
				if !m.Tiradas[cix] {
					info.NuestrasCartas[i/2] = append(
						info.NuestrasCartas[i/2],
						abs.Abstraer(c, &p.Ronda.Muestra),
					)
				}
			}
		}
		m = p.Ronda.GetSiguiente(*m)
	}
}

func (info *InfosetRondaBase) setManojos_en_juego(p *pdt.Partida, manojo *pdt.Manojo) {
	// 3. manojos que se fueron al mazo:

	// 	formato:
	// 		0.1.1.0--
	// done 1 significa que se fue al mazo

	// lo jugadores empezando por mi que se fueron
	info.ManojosEnJuego = make([]bool, len(p.Ronda.Manojos))
	m := manojo
	for i := 0; i < len(p.Ronda.Manojos); i++ {
		info.ManojosEnJuego[i] = !m.SeFueAlMazo
		m = p.Ronda.GetSiguiente(*m)
	}
}

func (info *InfosetRondaBase) setEnvido(p *pdt.Partida) {
	// 4. envido:
	info.Envido = p.Ronda.Envite.Estado.String()
}

func (info *InfosetRondaBase) setTruco(p *pdt.Partida) {
	// 5. el truco:
	info.Truco = p.Ronda.Truco.Estado.String()
}

func (info *InfosetRondaBase) setChi(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	chi_i pdt.A,
	abs abs.IAbstraccion,

) {
	// 6. chi(I):
	// 	las acciones que puede tomar son (e.g.):
	// 	x.y.z.1.0.1.1.0.1
	// 	done x,y,x \in \mathbb{N} se corresponde con la cantidad de cartas de
	// 	cada bucket de la abstraccion que puede tirar.
	// 	ejemplo usando la abstraccion A1:
	// 		`1.0.2` significa que puede jugar:
	// 			- 1 carta de tipo 0 (comunes)
	// 			- niguna carta de tipo 1 (mata)
	// 			- 2 cartas de tipo 2 (pieza)
	// 	que puede tomar (en el orden establecido)

	n := abs.Len()
	counter := make([]int, n) // tamano fijo (num de buckets)

	// como el bucket (el nivel de abstraccion) es el indice
	// entonces [0, 1, 0, 0, 2, 0]
	// significa que puede tirar 1 del bucket #1, y 2 del bucket #4
	// segun la escala de abs
	// en este caso, info.Chi_len() por lo tanto retornará 2.
	// ya que puede tirar cartas de hasta 2 buckets diferentes.

	// cartas
	for i := 0; i < 3; i++ {
		if cartaHabilitada := chi_i[i]; cartaHabilitada {
			c := manojo.Cartas[i]
			bucket := abs.Abstraer(c, &p.Ronda.Muestra)
			counter[bucket]++
		}
	}

	resto := make([]int, len(chi_i[3:])) // tamano fijo
	for ix, v := range chi_i[3:] {
		if v {
			resto[ix] = 1
		} else {
			resto[ix] = 0
		}
	}

	info.Chi = append(counter, resto...)
}

// varia la implementacion, de abstraccion en abstraccion
// al igual que setChi
func (info *InfosetRondaBase) ChiLen() int {
	// finalmente
	// cuento cuantos de estos buckets+acciones son "positivos":
	// esto lo uso para crear el RNode; no para el Infoset en sí
	chi_len := 0
	for _, a := range info.Chi {
		if a > 0 {
			chi_len += 1
		}
	}
	return chi_len
}

// Chi usa un vector de tamano variable que depende de la granularidad de la
// abstraccion usada. El simulador por su parte usa un array fijo de tamano 15:
// 3 cartas posibles (en orden) + 12 acciones.
// Para hacer que los infosets sean interoperables con el simulador se crea esta
// funcion.
// Transforma de un "Chi" a un `[]pdt.IJugada`
// PRE: ya se "inicio" el infoset
func (info *InfosetRondaBase) Iterable(

	p *pdt.Partida,
	m *pdt.Manojo,
	aixs pdt.A, // array de 15 acciones (bool): 3 cartas + 12 "jugadas"
	abs abs.IAbstraccion,

) []pdt.IJugada {

	// el slice a retornar sera como maximo de ese tamano
	res := make([]pdt.IJugada, 0, 3+12)

	// agrego primero las tiradas, luego las acciones
	// el resto de las acciones son faciles de "traducir"/adaptar
	n := abs.Len()

	// Otra vez lo de los buckets:
	// un slice de slices.
	// el primer indice indica el bucket;
	// luego se incluye un slice de indice absolutos a la posicion de la carta
	// que se puede tirar.
	// e.g.,
	// [[0,2], nil, nil, [1], nil]
	// significa que del bucket #0 puede tirar las cartas 0 y 2
	// y del bucket #3 puede tirar la carta 1
	counter := make([][]int, n)
	for i := 0; i < 3; i++ {
		noLaTiro := !m.Tiradas[i]
		laPuedeTirar := aixs[i]
		cartaHabilitada := noLaTiro && laPuedeTirar
		if cartaHabilitada {
			c := m.Cartas[i]
			bucket := abs.Abstraer(c, &p.Ronda.Muestra)
			if counter[bucket] == nil {
				counter[bucket] = []int{i}
			} else {
				counter[bucket] = append(counter[bucket], i)
			}
		}
	}

	// ahora, en base al slice anterior `counter` voy a agregar algunas jugadas
	// de tipo `pdt.TirarCarta` al retorno `res`.
	// si para un bucket dado hay mas de una carta que puedo tirar con esa
	// abstraccion, elijo solo una de ellas al azar:
	for _, bucketCount := range counter {
		if len(bucketCount) > 0 {
			// elijo una al azar *de ese bucket*
			randomIx := rand.Intn(len(bucketCount))
			cartaIx := bucketCount[randomIx]
			res = append(
				res,
				&pdt.TirarCarta{
					JID:   m.Jugador.ID,
					Carta: *m.Cartas[cartaIx],
				},
			)
		}
	}

	// se que las acciones empeizan luego de que termina el vector de tamano
	// variable determinado por la abstraccion:
	restoDeAcciones := info.Chi[n:]
	for ix, a := range restoDeAcciones {
		if a > 0 {
			canonicalAix := ix + 3
			// le sumo +3 ya que no estoy contando las 3 primeras acciones que son las
			// asociadas con las cartas;
			// como si fuera un offset
			res = append(
				res,
				pdt.ToJugada(
					p,
					p.Ronda.JIX(m.Jugador.ID),
					canonicalAix))
		}
	}

	return res
}

func (info *InfosetRondaBase) setResultadoManos(p *pdt.Partida, manojo *pdt.Manojo) {
	// 7. balance manos
	info.ResultadoManos = make([]string, 3)

	for mix := 0; mix < 3; mix++ {
		if int(p.Ronda.ManoEnJuego) <= mix {
			info.ResultadoManos[mix] = "?" // no se termino de definir
		} else {
			mano := &p.Ronda.Manos[mix]
			if mano.Resultado == pdt.Empardada {
				info.ResultadoManos[mix] = "parda"
			} else {
				laGanamos := p.Ronda.Manojo(mano.Ganador).Jugador.Equipo == manojo.Jugador.Equipo
				if laGanamos {
					info.ResultadoManos[mix] = "ganada"
				} else {
					info.ResultadoManos[mix] = "perdida"
				}
			}
		}
	}
}

func (info *InfosetRondaBase) setRonda(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	abs abs.IAbstraccion,

) {
	// 8. estado de la mano actual

	// De la mano actual, las cartas mas altas tirada por nosotros, por el
	// oponente y si vamos ganando o no.

	// formato:
	// 	c1.c2.x--
	// 	donde:
	// 		- c1 carta nuestra mas alta (string vacio si no)
	// 		- c2 carta oponente mas alta (string vacio si no)
	// 		- x = 0 si vamos perdiendo, 1 empatados o 2 ganando

	maxWe, maxOp, vamos := vamos(p, manojo, abs)

	estadoRonda := struct {
		Max_us int
		Max_op int
		Vamos  string
	}{
		maxWe,
		maxOp,
		vamos,
	}

	info.ManoActual = estadoRonda
}

func (info *InfosetRondaBase) Hash(h hash.Hash) string {
	// h := sha1.New()
	hsep := []byte(sep)

	// 1
	h.Write([]byte(strconv.Itoa(info.Muestra)))
	h.Write(hsep)

	// 2
	// paso de un array de abstracciones a un array de primos
	nuestrasCartas := make([]int, len(info.NuestrasCartas))
	for mix, manojo := range info.NuestrasCartas {
		manojoPrimeID := 1
		for _, abstraccion := range manojo {
			manojoPrimeID *= utils.AllPrimes[abstraccion]
		}
		nuestrasCartas[mix] = manojoPrimeID
	}
	bs, _ := json.Marshal(nuestrasCartas)
	h.Write([]byte(bs))
	h.Write(hsep)

	// 3
	bs, _ = json.Marshal(info.ManojosEnJuego)
	h.Write(bs)
	h.Write(hsep)

	// 4
	h.Write([]byte(info.Envido))
	h.Write(hsep)

	// 5
	h.Write([]byte(info.Truco))
	h.Write(hsep)

	// 6
	bs, _ = json.Marshal(info.ResultadoManos)
	h.Write([]byte(bs))
	h.Write(hsep)

	// 7
	manoActual := fmt.Sprintf("%d.%d.%s",
		info.ManoActual.Max_us, info.ManoActual.Max_op, info.ManoActual.Vamos)
	h.Write([]byte(manoActual))
	h.Write(hsep)

	// 8
	bs, _ = json.Marshal(info.Chi)
	h.Write(bs)
	// h.Write(hsep) // <- not necessary

	return hex.EncodeToString(h.Sum(nil))
}

func (info *InfosetRondaBase) Dump(indent bool) string {
	var bs []byte = nil
	if indent {
		bs, _ = json.MarshalIndent(info, "", "\t")
	} else {
		bs, _ = json.Marshal(info)
	}
	return string(bs)
}

func MkInfoset1(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	chi_i pdt.A,
	abs abs.IAbstraccion,

) Infoset {

	info := &InfosetRondaBase{
		Vision: manojo.Jugador.ID, // <- tiene motivos solo depurativos
	}

	info.setMuestra(p)
	info.setNuestras_Cartas(p, manojo, abs)
	info.setManojos_en_juego(p, manojo)
	info.setEnvido(p)
	info.setTruco(p)
	info.setChi(p, manojo, chi_i, abs)
	info.setResultadoManos(p, manojo)
	info.setRonda(p, manojo, abs)

	return info
}
