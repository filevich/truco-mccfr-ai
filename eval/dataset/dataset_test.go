package dataset

import (
	"os"
	"testing"

	"github.com/truquito/gotruco/pdt"
)

func TestEntryOverride2p(t *testing.T) {
	// Use os.Stat to get file info. It returns an error if the file doesn't exist.
	_, err := os.Stat("t1k22.json")

	// Check if an error occurred AND if that error indicates the file does not exist.
	if os.IsNotExist(err) {
		t.Skip("skipping test: t1k22 file not found")
		return
	}

	ds := LoadDataset("t1k22.json")

	p, _ := pdt.NuevaPartida(
		pdt.A20,
		[]string{"Alice"},
		[]string{"Bob"},
		2,
		true)

	if ok := len(ds) == 1_000 && len(ds[0]) == 79 &&
		ds[0][0].Muestra.Palo == pdt.Copa && ds[0][0].Muestra.Valor == 3; !ok {
		t.Error("los datos cargados no coinciden con los esperados")
	}

	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))

	// intercambio posiciones
	p.Swap()
	// reseteo
	p.Ronda.Reset(0)
	// sobre escribo los valores de las cartas
	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))
}

func TestEntryOverride4p(t *testing.T) {
	// Use os.Stat to get file info. It returns an error if the file doesn't exist.
	_, err := os.Stat("t1k22.json")

	// Check if an error occurred AND if that error indicates the file does not exist.
	if os.IsNotExist(err) {
		t.Skip("skipping test: t1k22 file not found")
		return
	}

	ds := LoadDataset("t1k22.json")

	p, _ := pdt.NuevaPartida(
		pdt.A20,
		[]string{"Alice", "Ariana"},
		[]string{"Bob", "Ben"},
		2,
		true)

	if ok := len(ds) == 1_000 && len(ds[0]) == 79 &&
		ds[0][0].Muestra.Palo == pdt.Copa && ds[0][0].Muestra.Valor == 3; !ok {
		t.Error("los datos cargados no coinciden con los esperados")
	}

	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))

	// intercambio posiciones
	p.Swap()
	// reseteo
	p.Ronda.Reset(0)
	// sobre escribo los valores de las cartas
	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))
}

func TestEntryOverride6p(t *testing.T) {
	// Use os.Stat to get file info. It returns an error if the file doesn't exist.
	_, err := os.Stat("t1k22.json")

	// Check if an error occurred AND if that error indicates the file does not exist.
	if os.IsNotExist(err) {
		t.Skip("skipping test: t1k22 file not found")
		return
	}

	ds := LoadDataset("t1k22.json")

	p, _ := pdt.NuevaPartida(
		pdt.A20,
		[]string{"Alice", "Ariana", "Anna"},
		[]string{"Bob", "Ben", "Bill"},
		2,
		true)

	if ok := len(ds) == 1_000 && len(ds[0]) == 79 &&
		ds[0][0].Muestra.Palo == pdt.Copa && ds[0][0].Muestra.Valor == 3; !ok {
		t.Error("los datos cargados no coinciden con los esperados")
	}

	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))

	// intercambio posiciones
	p.Swap()
	// reseteo
	p.Ronda.Reset(0)
	// sobre escribo los valores de las cartas
	ds[0][0].Override(p)

	if ok := p.Ronda.Muestra.Palo == pdt.Copa && p.Ronda.Muestra.Valor == 3; !ok {
		t.Error("los datos cargados de la muestra no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Bob").Cartas[0].Palo == pdt.Copa &&
		p.Ronda.Manojo("Bob").Cartas[0].Valor == 11; !ok {
		t.Error("los datos cargados en el manojo de Bob no coinciden con los esperados")
	}

	if ok := p.Ronda.Manojo("Alice").Cartas[0].Palo == pdt.Oro &&
		p.Ronda.Manojo("Alice").Cartas[0].Valor == 4; !ok {
		t.Error("los datos cargados en el manojo de Alice no coinciden con los esperados")
	}

	t.Log(pdt.Renderizar(p))
}
