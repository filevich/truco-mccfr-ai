package dataset

import (
	"encoding/json"
	"io"
	"os"

	"github.com/truquito/gotruco/pdt"
)

type Row struct {
	Muestra pdt.Carta        `json:"muestra"`
	Manojos [6][3]*pdt.Carta `json:"manojos"`
}

func (e *Row) Override(p *pdt.Partida) {
	p.Ronda.Muestra = e.Muestra
	for mix := range p.Ronda.Manojos {
		p.Ronda.Manojos[mix].Cartas = e.Manojos[mix]
	}
	p.Ronda.CachearFlores(true)
}

type Dataset [][]*Row

func LoadDataset(filepath string) Dataset {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bs, _ := io.ReadAll(file)

	var ds [][]*Row
	json.Unmarshal(bs, &ds)

	return ds
}
