package abs

import (
	"encoding/json"
	"strings"

	"github.com/truquito/gotruco/pdt"
)

// Cards abstractions:
// A1, A2, A3, B

type IAbstraction interface {
	Len() int                                      // returns abstraction's number of buckets
	Abstract(c *pdt.Carta, muestra *pdt.Carta) int // returns card's bucket
	String() string                                // retrns abstraction's ID
	MarshalJSON() ([]byte, error)
}

type Abstractor_ID string

const (
	A1_ID   Abstractor_ID = "a1"
	B_ID    Abstractor_ID = "b"
	A2_ID   Abstractor_ID = "a2"
	A3_ID   Abstractor_ID = "a3"
	NULL_ID Abstractor_ID = "null"
)

func ParseAbs(aID string) IAbstraction {
	switch Abstractor_ID(strings.ToLower(aID)) {
	case A1_ID:
		return &A1{}
	case B_ID:
		return &B{}
	case A2_ID:
		return &A2{}
	case A3_ID:
		return &A3{}
	case NULL_ID:
		return &Null{}
	}

	panic("abstraction unknown")
}

type A1 struct{}

// abstraction A1:
//
// {2,4,5,11,10} (de la muestra)
// -------
// {1,1,7,7} (matas)
// -------
// {3,2,1,12,11,10,7,6,5,4} (resto)

func (a A1) String() string {
	return string(A1_ID)
}

func (a A1) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a A1) Len() int {
	return 3 // <--- its max bucket +1
}

func (a A1) Abstract(c *pdt.Carta, muestra *pdt.Carta) int {
	if c.EsPieza(*muestra) {
		return 2
	} else if c.EsMata() {
		return 1
	}

	return 0
}

/*






 */

type A2 struct{}

// abstraction A2:
//
// {2,4,5} (de la muestra)
// {11,10} (de la muestra)
// -------
// {1,1} (matas)
// {7,7}
// -------
// {3,2,1}
// {12, 11, 10}
// {7,6,5,4}

func (a A2) String() string {
	return string(A2_ID)
}

func (a A2) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a A2) Len() int {
	return 7 // <--- its max bucket +1
}

func (a A2) Abstract(c *pdt.Carta, muestra *pdt.Carta) int {
	if c.EsPieza(*muestra) {
		if c.Valor == 2 || c.Valor == 4 || c.Valor == 5 {
			return 6
		} else {
			return 5
		}
	} else if c.EsMata() {
		if c.Valor == 1 {
			return 4
		} else {
			return 3
		}
	} else {
		if c.Valor <= 3 {
			return 2
		} else if 10 <= c.Valor && c.Valor <= 12 {
			return 1
		} else {
			return 0
		}
	}
}

type A3 struct{}

// abstraction A3:

func (a A3) String() string {
	return string(A3_ID)
}

func (a A3) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a A3) Len() int {
	return 19
}

var A3_map = map[int]int{
	34: 18, // 2 de la muestra
	33: 17, // 4 de la muestra
	32: 16, // 5 de la muestra
	31: 15, // 11 de la muestra
	30: 14, // 10 de la muestra
	23: 13, // 1 espada
	22: 12, // 1 basto
	21: 11, // 7 espada
	20: 10, // 7 oro
	19: 9,  // 3 ?
	18: 8,  // 2 ?
	17: 7,  // 1 ?
	16: 6,  // 12 ?
	15: 5,  // 11 ?
	14: 4,  // 10 ?
	13: 3,  // 7 ?
	12: 2,  // 6 ?
	11: 1,  // 5 ?
	10: 0,  // 4 ?
}

func (a A3) Abstract(c *pdt.Carta, muestra *pdt.Carta) int {
	return A3_map[c.CalcPoder(*muestra)]
}

/*






 */

type B struct{}

// abstraction A1:
//
// {2,4,5,11,10} (de la muestra)
// -------
// {1,1,7,7} (matas)
// -------
// {3,2,1,12,11,10}
// -------
// {7,6,5,4} (resto)

func (a B) String() string {
	return string(B_ID)
}

func (a B) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a B) Len() int {
	return 4 // <--- its max bucket +1
}

func (a B) Abstract(c *pdt.Carta, muestra *pdt.Carta) int {
	if c.EsPieza(*muestra) {
		return 3
	} else if c.EsMata() {
		return 2
	} else if c.Valor <= 3 || c.Valor >= 10 {
		return 1
	}

	return 0
}

/*






 */

type Null struct{}

// abstraction Null:

/*

si estamos jugando al truco, y en una ronda la muestra es el 3 de oro,
entonces, tanto el 4 de bastro como el 4 de espada determinan infosets
diferentes (a pesar de tener el mismo "poder").
Los queremos distinguir.
Cada carta tiene asociada un único número.

*/

/*
 *  Barajas; orden absoluto:
 *  ----------------------------------------------------------
 * | ID	| Carta	    ID | Carta	  ID | Carta	    ID | Carta |
 * |---------------------------------------------------------|
 * | 00 | 1,basto   10 | 1,copa   20 | 1,espada   30 | 1,oro |
 * | 01 | 2,basto   11 | 2,copa   21 | 2,espada   31 | 2,oro |
 * | 02 | 3,basto   12 | 3,copa   22 | 3,espada   32 | 3,oro |
 * | 03 | 4,basto   13 | 4,copa   23 | 4,espada   33 | 4,oro |
 * | 04 | 5,basto   14 | 5,copa   24 | 5,espada   34 | 5,oro |
 * | 05 | 6,basto   15 | 6,copa   25 | 6,espada   35 | 6,oro |
 * | 06 | 7,basto   16 | 7,copa   26 | 7,espada   36 | 7,oro |
 *  ----------------------------------------------------------
 * | 07 |10,basto   17 |10,copa   27 |10,espada   37 |10,oro |
 * | 08 |11,basto   18 |11,copa   28 |11,espada   38 |11,oro |
 * | 09 |12,basto   19 |12,copa   29 |12,espada   39 |12,oro |
 *  ----------------------------------------------------------
 */

func (a Null) String() string {
	return string(NULL_ID)
}

func (a Null) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a Null) Len() int {
	return 40
}

// notar que esta abstaccion es independiente de la muestra
func (a Null) Abstract(c *pdt.Carta, muestra *pdt.Carta) int {
	return int(c.ID())
}
