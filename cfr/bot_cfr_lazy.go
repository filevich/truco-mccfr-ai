package cfr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

// on avg, it takes 6.2s to play an entire game of Truco between Simple bot
// and BotLazyCFR, on a 512GiB M2 MBA.
// if we play 2*1,000 matches then that implies a total of 3,44 hours :(

func ReadLimit(
	hash string,
	limit int64,
	r io.Reader,
) (
	found bool,
	line string,

) {
	scanner := bufio.NewScanner(r)

	// buff size
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var totalRead int64 = 0

	for scanner.Scan() {
		bytesRead := len(scanner.Bytes())
		totalRead += int64(bytesRead)
		line = scanner.Text()
		found = strings.HasPrefix(line, hash)
		if found || totalRead >= limit {
			break
		}
	}

	return found, line
}

func ReadLimitMultithread(
	hash string,
	limit int64,
	r io.Reader,
	done chan struct{},
) (
	found bool,
	line string,

) {
	scanner := bufio.NewScanner(r)

	// buff size
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var totalRead int64 = 0

	for scanner.Scan() {
		select {
		case <-done:
			// Terminate early if done signal received
			return false, ""
		default:
			bytesRead := len(scanner.Bytes())
			totalRead += int64(bytesRead)
			line = scanner.Text()
			found = strings.HasPrefix(line, hash)
			if found {
				return found, line
			} else if totalRead >= limit {
				return false, ""
			}
		}
	}

	return false, ""
}

type BotLazyCFR struct {
	ID       string
	Filepath string
	trainer  ITrainer

	Threads  int64
	fs       []*os.File
	fileSize int64
}

func (b *BotLazyCFR) Initialize() {
	log.Println("initing lazy")

	if b.fs != nil {
		return
	}

	b.fs = make([]*os.File, b.Threads)
	for i := 0; i < int(b.Threads); i++ {
		f, _ := os.Open(b.Filepath)
		b.fs[i] = f
	}

	stat, err := b.fs[0].Stat()
	if err != nil {
		panic(err)
	}
	b.fileSize = stat.Size()

	b.trainer = LoadModel(b.Filepath, false, 1_000_000, true)
	log.Println("done lazy loading", b.trainer.GetAbs().String())
}

func (b *BotLazyCFR) Free() {
	for i := 0; i < int(b.Threads); i++ {
		b.fs[i].Close()
	}
}

func (b *BotLazyCFR) UID() string {
	return b.ID
}

func (b *BotLazyCFR) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *BotLazyCFR) ResetCatch() {}

func (b *BotLazyCFR) _resetfilePtr() {
	blockSize := b.fileSize / b.Threads
	for i := 0; i < int(b.Threads); i++ {
		b.fs[i].Seek(int64(i)*blockSize, io.SeekStart)
	}
}

func _branch(s, match string) (head, tail string) {
	branchArray := strings.SplitN(s, match, 2)
	return branchArray[0], branchArray[1]
}

func (b *BotLazyCFR) Find(hash string) (rnode *RNode, err error) {
	wg := sync.WaitGroup{}
	resultChan := make(chan string, 1) // Buffer for one result
	done := make(chan struct{})

	// jumps
	b._resetfilePtr()

	limit := b.fileSize / b.Threads

	for i := 0; i < int(b.Threads); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if found, line := ReadLimitMultithread(hash, limit, b.fs[i], done); found {
				resultChan <- line // Send true if found
				close(done)        // Signal other goroutines to stop
			}
		}(i)
	}

	wg.Wait() // Ensure goroutines finish before closing channel
	close(resultChan)

	if l := <-resultChan; len(l) > 0 {
		_, tail := _branch(l, hash+" ")
		rnode = &RNode{}
		if err := json.Unmarshal([]byte(tail), rnode); err != nil {
			return nil, fmt.Errorf("no se pudo parsear la linea %s", hash)
		}
		return rnode, nil
	}

	return nil, fmt.Errorf("hash not found %s", hash)
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
	// log.Println(b.trainer, inGameID, active_player, p)
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
