package cfr_test

import (
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/filevich/truco-mccfr-ai/cfr"
)

var (
	sampleModel = "/Users/jp/Downloads/cluster/train-cfr/models/2p/irb-a3/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model"
	sampleHash  = "1207b37ae629fa3d2cb8aa11bbc56602b2a0e389"
)

func TestLinealRead(t *testing.T) {
	filename := sampleModel
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stat, _ := f.Stat()
	fileSize := stat.Size()

	hash := sampleHash

	{
		limit := int64(937_619_324)
		start := time.Now()
		found, line := cfr.ReadLimit(hash, limit, f)
		if ok := found; !ok {
			t.Error("it should find the line for the given limit")
			t.Fail()
		}
		t.Log("delta", time.Since(start))
		t.Log("fileSize", fileSize)
		t.Log("found", found)
		if found {
			t.Log("line", line)
		}
	}

	{
		limit := int64(619_324)
		found, _ := cfr.ReadLimit(hash, limit, f)
		if ok := !found; !ok {
			t.Error("it should NOT find the line for the given limit")
			t.Fail()
		}
	}

}

func TestMultithreadRead(t *testing.T) {
	wg := sync.WaitGroup{}
	resultChan := make(chan string, 1) // Buffer for one result
	done := make(chan struct{})

	filename := sampleModel
	threads := int64(6)

	fs := make([]*os.File, threads)
	for i := 0; i < int(threads); i++ {
		f, _ := os.Open(filename)
		fs[i] = f
	}

	defer func() {
		for i := 0; i < int(threads); i++ {
			fs[i].Close()
		}
	}()

	stat, _ := fs[0].Stat()
	fileSize := stat.Size()

	// jumps
	blockSize := fileSize / threads

	for i := 0; i < int(threads); i++ {
		fs[i].Seek(int64(i)*blockSize, io.SeekStart)
	}

	hash := sampleHash
	// hash := "be3445aa21c7798c22f531906962185d09dccd1f"
	limit := blockSize

	for i := 0; i < int(threads); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if found, line := cfr.ReadLimitMultithread(hash, limit, fs[i], done); found {
				resultChan <- line // Send true if found
				close(done)        // Signal other goroutines to stop
			}
		}(i)
	}

	wg.Wait() // Ensure goroutines finish before closing channel
	close(resultChan)

	if s := <-resultChan; len(s) > 0 {
		fmt.Println("hash found!", s)
	} else {
		fmt.Println("hash not found :(")
	}
}

func TestCFRLazy1(t *testing.T) {
	b := &cfr.BotLazyCFR{
		ID:       "example",
		Filepath: sampleModel,
	}
	b.Initialize()
	node, err := b.Find(sampleHash)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("node", node)
}
