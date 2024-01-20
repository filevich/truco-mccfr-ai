package info

import (
	"github.com/filevich/truco-cfr/abs"
	"github.com/truquito/truco/pdt"
)

// type Infoset interface {
// 	Hash() string
// 	Chi_len() int
// 	Dump(indent bool) string
// 	Iterable(
// 		p *pdt.Partida,
// 		m *pdt.Manojo,
// 		aixs pdt.A,
// 		abs abs.IAbstraccion,
// 	) []pdt.IJugada
// }

// S, SE, C, R, plus pro max ultra

//                            ╔════════════════╗
// ┌4─┐          ┌4─┐12┐6─┐   │ #Mano: Primera │
// │//│          │es│co│ba│   ╠────────────────╣
// └──┘          └──┘──┘──┘   │ Mano: Alice    │
//                            ╠────────────────╣
//    Ben          Ariana     │ Turno: Alice   │
// ╔════════════════════════╗ ╠────────────────╣
// ║                        ║ │ Puntuacion: 20 │
// ║ ┌4─┐2─┐                ║ ╚════════════════╝
// ║ │co│co│   ┌5─┐         ║ ╔════════════════╗
// ║ └──┘──┘   │es│         ║ │ Envite: noCant │
// ║           └──┘         ║ │ Por:           │
// ║                        ║ ╠────────────────╣
// ║                        ║ │ Truco: noGritad│
// ╚════════════════════════╝ │ Por:           │
//   Alice          Bob       ╚════════════════╝
// 	   ↑             ❀        ╔═══════╦╦═══════╗
// ┌10┐12┐10┐    ┌6─┐3─┐11┐   │ Alice ││ Bob   │
// │ba│ba│or│    │//│//│//│   ╠───────┼┼───────╣
// └──┘──┘──┘    └──┘──┘──┘   │   0   ││   0   │
//                            ╚═══════╩╩═══════╝

// guarda toda la info se la ronda
// pero no guarda:
//   - ninguna informacion relacionada a la partida. I.e., puntaje y puntuacion
//   - "el historial" de los envites/truco
//   - "el historial" de jugadas. e.g., Bob:Mazo Ben:1C ~ Ben:1C Bob:Mazo

type InfosetRondaPro struct {
	// `muestra`
	// la almacena de forma pura; sin abstracción.
	// no es posible abstraer la muestra porque la función de abstraer depende
	// de la muestra misma.
	// notar que almacenar solor el valor de la muestra (i.e., el número)
	// también es un tipo de abstracción.
	muestra int

	// `num_mano_actual`: int
	numMano int

	// `rixMe` ~ RIX: who?
	rixMe int

	// `turno` ~ RIX who?
	rixTurno int

	// `ManojosEnJuego` quiénes se fueron al mazo y quiénes siguen en pie?
	manojosEnJuego []bool

	// `nuestrasCartas` representa nuestras cartas.
	// Por cada jugador/manojo de nuestro equipo se almacena un único entero.
	// Estos enteros se calculan como la multiplicación de todos los primos
	// asociados a cada carta del manojo.
	// De esta forma, cartas diferentes generan identificadores diferentes
	// pero el orden (permutación) no altera el resultado.
	// El indice 0 se corresponde con las cartas del manojo más cercano de nuestro
	// equipo al Mano de la ronda actual.
	nuestrasCartas []int

	// tiradas: [(abs(carta),who?)]
	tiradasCartas []int
	tiradasWho    []int

	// envite_history: [(estado,who?)] <- inc says

	// truco_history: [(estado,who?)] <- inc says

	// flores: [who?]

}

func (info *InfosetRondaPro) setMuestra(p *pdt.Partida) {
	info.muestra = int(p.Ronda.Muestra.ID())
}

func (info *InfosetRondaPro) setNumMano(p *pdt.Partida) {
	info.numMano = int(p.Ronda.ManoEnJuego)
}

func (info *InfosetRondaPro) setMano(p *pdt.Partida, m *pdt.Manojo) {
	info.rixMe = RIX(p, m)
}

func (info *InfosetRondaPro) setTurno(p *pdt.Partida) {
	info.rixTurno = RIX(p, p.Ronda.GetElTurno())
}

func (info *InfosetRondaPro) setManojosEnJuego(p *pdt.Partida) {
	// tengo que empezar a iterar a partir del JIX del MANO
	n := len(p.Ronda.Manojos)
	info.manojosEnJuego = make([]bool, n)
	m := p.Ronda.GetElMano()
	for i := 0; i < n; i++ {
		info.manojosEnJuego[i] = !m.SeFueAlMazo
		m = p.Ronda.GetSiguiente(*m)
	}
}

func (info *InfosetRondaPro) setNuestrasCartas(

	p *pdt.Partida,
	e pdt.Equipo,
	a abs.IAbstraccion,

) {
	// cada equipo tiene n/2 manojos donde n es la cantidad de jugadores en la
	// partida.
	// Nota: dividir entre 2 es igual a hacer un shift right `n >> 1`
	n := len(p.Ronda.Manojos)
	info.nuestrasCartas = make([]int, n>>1)

	// tengo que empezar a iterar a partir del JIX del MANO
	m := p.Ronda.GetElMano()
	for i := 0; i < n; i++ {
		if esDeNuestroEquipo := m.Jugador.Equipo == e; esDeNuestroEquipo {
			info.nuestrasCartas[i] = PrimifyManojo(m, &p.Ronda.Muestra, a)
		}
		m = p.Ronda.GetSiguiente(*m)
	}
}

func (info *InfosetRondaPro) setTiradas(

	p *pdt.Partida,
	e pdt.Equipo,
	a abs.IAbstraccion,

) {
	// cada equipo tiene n/2 manojos donde n es la cantidad de jugadores en la
	// partida.
	// Nota: dividir entre 2 es igual a hacer un shift right `n >> 1`
	n := len(p.Ronda.Manojos)
	info.nuestrasCartas = make([]int, n>>1)

	// tengo que empezar a iterar a partir del JIX del MANO
	m := p.Ronda.GetElMano()
	for i := 0; i < n; i++ {
		if esDeNuestroEquipo := m.Jugador.Equipo == e; esDeNuestroEquipo {
			info.nuestrasCartas[i] = PrimifyManojo(m, &p.Ronda.Muestra, a)
		}
		m = p.Ronda.GetSiguiente(*m)
	}
}

func NewInfosetUltra(

	p *pdt.Partida,
	manojo *pdt.Manojo,
	abs abs.IAbstraccion,

) Infoset {

	info := &Infoset1{
		Vision: manojo.Jugador.ID, // <- tiene motivos solo depurativos
	}

	info.setMuestra(p)
	info.setNuestras_Cartas(p, manojo, abs)
	info.setManojos_en_juego(p, manojo)
	info.setEnvido(p)
	info.setTruco(p)
	// info.setChi(p, manojo, chi_i, abs)
	info.setResultadoManos(p, manojo)
	info.setRonda(p, manojo, abs)

	return info
}
