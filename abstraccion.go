package trucocfr

import (
	"encoding/json"

	"github.com/truquito/truco/pdt"
)

// Cards abstractions:
// A1, A2, A3, B1

type IAbstraccion interface {
	Len() int                                      // returns abstraction's number of buckets
	Abstraer(c *pdt.Carta, muestra *pdt.Carta) int // returns card's bucket
	String() string                                // retrns abstraction's ID
	MarshalJSON() ([]byte, error)
}

type Abstractor_ID string

const (
	A1_ID Abstractor_ID = "a1"
	A2_ID Abstractor_ID = "a2"
	A3_ID Abstractor_ID = "a3"
	B1_ID Abstractor_ID = "b1"
)

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

func (a A1) Abstraer(c *pdt.Carta, muestra *pdt.Carta) int {
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

func (a A2) Abstraer(c *pdt.Carta, muestra *pdt.Carta) int {
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

func (a A3) Abstraer(c *pdt.Carta, muestra *pdt.Carta) int {
	return A3_map[c.CalcPoder(*muestra)]
}

/*






 */

type B1 struct{}

// abstraction A1:
//
// {2,4,5,11,10} (de la muestra)
// -------
// {1,1,7,7} (matas)
// -------
// {3,2,1,12,11,10}
// -------
// {7,6,5,4} (resto)

func (a B1) String() string {
	return string(B1_ID)
}

func (a B1) MarshalJSON() ([]byte, error) {
	str := a.String()
	return json.Marshal(str)
}

func (a B1) Len() int {
	return 4 // <--- its max bucket +1
}

func (a B1) Abstraer(c *pdt.Carta, muestra *pdt.Carta) int {
	if c.EsPieza(*muestra) {
		return 3
	} else if c.EsMata() {
		return 2
	} else if c.Valor <= 3 || c.Valor >= 10 {
		return 1
	}

	return 0
}
