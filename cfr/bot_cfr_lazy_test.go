package cfr_test

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/filevich/truco-ai/cfr"
)

var (
	sampleModel = "/Users/jp/Downloads/cluster/train-cfr/models/2p/irb-a3/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model"
	sampleHash  = "1207b37ae629fa3d2cb8aa11bbc56602b2a0e389"
)

func readLimit(
	hash string,
	limit int64,
	r io.Reader,
) (
	found bool,
	line string,

) {
	scanner := bufio.NewScanner(r)

	// buff size
	// const maxCapacity = 1024 * 1024
	// buf := make([]byte, maxCapacity)
	// scanner.Buffer(buf, maxCapacity)

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
		found, line := readLimit(hash, limit, f)
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
		found, _ := readLimit(hash, limit, f)
		if ok := !found; !ok {
			t.Error("it should NOT find the line for the given limit")
			t.Fail()
		}
	}

}

func TestMultithreadRead(t *testing.T) {
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
	limit := blockSize

	for i := 0; i < int(threads); i++ {
		start := time.Now()
		found, _ := readLimit(hash, limit, fs[i])
		t.Log("delta", time.Since(start))
		t.Log("fileSize", fileSize)
		t.Log("found", found)
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
