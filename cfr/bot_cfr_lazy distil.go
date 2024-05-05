package cfr

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"

	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type BotLazyDistilCFR struct {
	ID       string
	Filepath string
	trainer  ITrainer

	Threads  int64
	fs       []*os.File
	fileSize int64

	// stats
	Hit  int
	Miss int
}

func (b *BotLazyDistilCFR) Initialize() {
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

func (b *BotLazyDistilCFR) Free() {
	for i := 0; i < int(b.Threads); i++ {
		b.fs[i].Close()
	}
	r := float32(b.Hit) / float32(b.Hit+b.Miss)
	log.Printf("hit=%d miss=%d hit_ratio=%.2f", b.Hit, b.Miss, r)
}

func (b *BotLazyDistilCFR) UID() string {
	return b.ID
}

func (b *BotLazyDistilCFR) Catch(*pdt.Partida, []enco.Envelope) {}

func (b *BotLazyDistilCFR) ResetCatch() {}

func (b *BotLazyDistilCFR) _resetfilePtr() {
	blockSize := b.fileSize / b.Threads
	for i := 0; i < int(b.Threads); i++ {
		b.fs[i].Seek(int64(i)*blockSize, io.SeekStart)
	}
}

func (b *BotLazyDistilCFR) Find(hash string) (a int, err error) {
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
		a, err := strconv.Atoi(tail)
		if err != nil {
			panic(err)
		}
		return a, nil
	}

	return -1, fmt.Errorf("hash not found %s", hash)
}

func (b *BotLazyDistilCFR) Action(

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
	aix, err := b.Find(hash)
	if err != nil {
		b.Miss++
		// log.Println(err)
		// 	panic(err)
	} else {
		b.Hit++
	}

	// obtengo el chi
	Chi := i.Iterable(p, active_player, aixs, b.trainer.GetAbs())

	if err != nil {
		aix = rand.Intn(len(Chi))
	}

	return Chi[aix], 0
}
