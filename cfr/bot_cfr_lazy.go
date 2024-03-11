package cfr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/truco/enco"
	"github.com/truquito/truco/pdt"
)

type BotLazyCFR struct {
	ID       string
	Filepath string
	trainer  ITrainer
	filePtr  *os.File
}

func (b *BotLazyCFR) Initialize() {
	fmt.Println("initing lazy")
	// lo cargo SOLO si no fue cargado aun
	if b.filePtr == nil {
		f, err := os.Open(b.Filepath)
		if err != nil {
			panic(err)
		}
		b.filePtr = f
		// b.
		b.trainer = LoadModel(b.Filepath, false, 1_000_000, true)
		fmt.Println("done lazy loading", b.trainer.GetAbs().String())
	}
}

func (b *BotLazyCFR) Free() {
	if b.filePtr != nil {
		b.filePtr.Close()
	}
}

func (b *BotLazyCFR) UID() string {
	return b.ID
}

func (b *BotLazyCFR) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *BotLazyCFR) ResetCatch() {}

func (b *BotLazyCFR) _resetfilePtr() {
	// call the Seek method first
	_, err := b.filePtr.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
}

func _branch(s, match string) (head, tail string) {
	branchArray := strings.SplitN(s, match, 2)
	return branchArray[0], branchArray[1]
}

func (b *BotLazyCFR) Find(hash string) (rnode *RNode, err error) {
	defer b._resetfilePtr()

	c, l, found := 0, "", false

	scanner := bufio.NewScanner(b.filePtr)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		c++
		l = scanner.Text()
		if found = strings.HasPrefix(l, hash); found {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("no se pudo parsear la linea %s", hash)
	}

	_, tail := _branch(l, hash+" ")

	rnode = &RNode{}
	if err := json.Unmarshal([]byte(tail), rnode); err != nil {
		return nil, fmt.Errorf("no se pudo parsear la linea %s", hash)
	}

	return rnode, nil
}

func (b *BotLazyCFR) Action(

	p *pdt.Partida,
	inGameID string,

) (

	pdt.IJugada,
	float32,

) {

	// pseudo jugador activo
	active_player := p.Manojo(inGameID)

	// obtengo el infoset
	aixs := pdt.GetA(p, active_player)
	// fmt.Println(b.trainer, inGameID, active_player, p)
	i := b.trainer.GetBuilder().Info(p, active_player, nil)
	hash, _ := i.Hash(b.trainer.GetBuilder().Hash), i.ChiLen()

	// creo un Rnode
	rnode, err := b.Find(hash)
	if err != nil {
		panic(err)
	}

	// obtengo la strategy
	strategy := rnode.GetAverageStrategy()
	aix := utils.Sample(strategy)

	// obtengo el chi
	Chi := i.Iterable(p, active_player, aixs, b.trainer.GetAbs())

	return Chi[aix], strategy[aix]
}
