package eval

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/truquito/truco/pdt"
)

type Entrada struct {
	Muestra pdt.Carta        `json:"muestra"`
	Manojos [6][3]*pdt.Carta `json:"manojos"`
}

func (e *Entrada) Override(p *pdt.Partida) {
	p.Ronda.Muestra = e.Muestra
	for mix := range p.Ronda.Manojos {
		p.Ronda.Manojos[mix].Cartas = e.Manojos[mix]
	}
	p.Ronda.CachearFlores(true)
}

type Dataset [][]*Entrada

func Load_dataset(filepath string) Dataset {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bs, _ := ioutil.ReadAll(file)

	var ds [][]*Entrada
	json.Unmarshal(bs, &ds)

	return ds
}
