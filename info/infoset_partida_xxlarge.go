package info

import (
	"encoding/json"
	"hash"
	"strconv"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetPartidaXXLarge struct {
	InfosetRondaXXLarge
	puntajeMe int
	puntajeOp int
}

func (info *InfosetPartidaXXLarge) setPuntajes(p *pdt.Partida, m *pdt.Manojo) {
	if somosAzules := m.Jugador.Equipo == pdt.Azul; somosAzules {
		info.puntajeMe = p.Puntajes[pdt.Azul]
		info.puntajeOp = p.Puntajes[pdt.Rojo]
	} else {
		info.puntajeMe = p.Puntajes[pdt.Rojo]
		info.puntajeOp = p.Puntajes[pdt.Azul]
	}
}

func (info *InfosetPartidaXXLarge) HashBytes(h hash.Hash) []byte {
	h.Reset()
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

	// 10. puntajes
	h.Write([]byte(strconv.Itoa(info.puntajeMe)))
	h.Write(hsep)
	h.Write([]byte(strconv.Itoa(info.puntajeOp)))

	return h.Sum(nil)
}

func infosetPartidaXXLargeFactory(

	a abs.IAbstraction,

) InfosetBuilder {

	return func(

		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,

	) Infoset {
		info := &InfosetPartidaXXLarge{}
		info.setPuntajes(p, m)
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

}
