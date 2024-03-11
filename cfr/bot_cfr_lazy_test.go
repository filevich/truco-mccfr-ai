package cfr_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/filevich/truco-ai/cfr"
)

var sample = "/Users/jp/Downloads/cluster/train-cfr/models/2p/irb-a3/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model"

func TestLazyRead(t *testing.T) {
	filename := sample
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	// scanner := bufio.NewScanner(f)
	c, l, found := 0, "", false
	hash := "1207b37ae629fa3d2cb8aa11bbc56602b2a0e389"

	scanner := bufio.NewScanner(f)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		c++
		l = scanner.Text()
		found = strings.HasPrefix(l, hash)
		if found {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	t.Log(c, found, l)
}

func TestCFRLazy1(t *testing.T) {
	b := &cfr.BotLazyCFR{
		ID:       "example",
		Filepath: sample,
	}
	b.Initialize()
	node, err := b.Find("4e2ead8dc8e4ae11c4c60e2271cfc6ba0b7f673b")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("node", node)
}
