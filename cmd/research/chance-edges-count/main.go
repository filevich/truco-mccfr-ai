package main

import (
	"log"

	"github.com/filevich/combinatronics"
	"github.com/filevich/truco-ai/utils"
)

func main() {

	// n := 40 // full deck
	n := 12
	// n := 1 + 3 + 3 // <- min 2p

	miniIDs := make([]int, n)
	for i := range miniIDs {
		miniIDs[i] = i
	}

	// mazo de cartas primas
	// miniIDs := []int{10, 20, 11, 21, 14, 24, 15, 25, 16, 26, 17, 27, 18, 28}
	// miniIDs := []int{20, 0, 26, 36, 12, 16, 5}

	var total uint64 = 0

	// todas las muestras posibles
	for _, muestraID := range miniIDs {
		resto := utils.CopyWithoutThese(miniIDs, muestraID)
		// todos mis manojos posibles
		for _, miManojoIDs := range combinatronics.Combs(resto, 3) {
			resto2 := utils.CopyWithoutThese(resto, miManojoIDs...)
			// todos sus manojos posibles
			total += uint64(len(combinatronics.Combs(resto2, 3)))
		}
	}

	log.Println("number of possible chance node edges:", total)
}
