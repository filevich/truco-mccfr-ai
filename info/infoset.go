package info

import (
	"hash"

	"github.com/filevich/truco-cfr/abs"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

const (
	sep = "--" // separador
	div = "."  // divisor
)

/*

notas:

  1. Los infosets deben computar el mismo hash independientemente del nombre de
     los jugadores. Por ende, deben ser computados en funcion de su posicion
     relativa a quien es "el mano actual". Por que? El indice de los jugadores
     es un artefacto puramente de implementacion; que en la vida real se asemeja
     con "la posicion en la mesa" Para un momento dado cualquiera en una partida
     de truco, el hecho de que los jugadores roten en +1 su posicion en la mesa
     no deberia de afectar el computo de los infosets. Lo que si deberia de
     afectar el computo de los infosets es si cambia su distancia a "el mano
     actual"; ya que en ese caso, la posicion estrategica (y no fisica) del
     jugador cambio.

  2. El vector de accion (Chi) tiene las opciones "1era" "2da" y "3era" segun el
     orden en que fueron repartidas por la Naturaleza. (Independitemente del
     valor.) El problema es que el infoset es agnostico a la permutacion y solo
     depende de la combinacion de los "poderes" de las cartas.

  3. De por si, con 4 o mas jugadores ya existe una abstraccion de las cartas
     chicas. Si jugamos con señas entonces "esto es nativo del truco".

*/

// Modelar a los infosets como una interfaz permite multiples implementaciones
// donde algunas pueden ser mas granulares que otras.
// Esto permite generar agentes mas o menos pesados segun se desee y asi
// administrar mejor los recursos.

type Infoset interface {
	HashBytes(hash.Hash) []byte
	Hash(hash.Hash) string
	ChiLen() int
	Dump(indent bool) string
	Iterable(
		p *pdt.Partida,
		m *pdt.Manojo,
		aixs pdt.A,
		abs abs.IAbstraction,
	) []pdt.IJugada
}

type InfosetBuilder func(

	p *pdt.Partida,
	m *pdt.Manojo,
	a abs.IAbstraction,
	msgs []enco.IMessage,

) Infoset

func ParseInfosetBuilder(ID string) InfosetBuilder {
	builders := map[string]InfosetBuilder{
		"InfosetRondaBase":  NewInfosetRondaBase,
		"InfosetRondaLarge": NewInfosetRondaLarge,
	}
	if maker, ok := builders[ID]; ok {
		return maker
	}

	panic("infoset impl. does not exists")
}
