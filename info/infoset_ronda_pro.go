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
// 							  ╠────────────────╣
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
// 							  ╚═══════╩╩═══════╝

// guarda toda la info se la ronda
// pero no guarda:
//   - ninguna informacion relacionada a la partida. I.e., puntaje y puntuacion
//   - "el historial" de los envites/truco
//   - "el historial" de jugadas. e.g., Bob:Mazo Ben:1C ~ Ben:1C Bob:Mazo

type InfosetRondaPro struct {
	// muestra
	// la almacena de forma pura; sin abstracción.
	// no es posible abstraer la muestra porque la función de abstraer depende
	// de la muestra misma.
	// notar que almacenar solor el valor de la muestra (i.e., el número)
	// también es un tipo de abstracción.
	muestra int

	// num_mano_actual: int
	manoActual int

	// mano: who?

	// turno: who?
	// envite_history: [(estado,who?)] <- inc says
	// truco_history: [(estado,who?)] <- inc says
	// flores: [who?]
	// nuestras: cartas
	// tiradas: [(abs(carta),who?)]
}

func (info *InfosetRondaPro) setMuestra(p *pdt.Partida) {
	info.muestra = int(p.Ronda.Muestra.ID())
}

func (info *InfosetRondaPro) setManoActual(p *pdt.Partida) {
	info.manoActual = int(p.Ronda.ManoEnJuego)
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
