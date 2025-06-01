package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/filevich/truco-mccfr-ai/bot"
	"github.com/filevich/truco-mccfr-ai/eval"
	"github.com/filevich/truco-mccfr-ai/eval/dataset"
	"github.com/filevich/truco-mccfr-ai/utils"
)

var (
	numPlayers              = 0
	agents     []eval.Agent = nil
)

func init() {
	// arg 0 - whole program (+1)
	// arg 1 - numPlayers:int (+1)
	// args 2..n - model|simple|random (at least 2 so, +2)

	// e.g.,
	// go run cmd/eval/bin/main.go 2 simple random /models/2p/a1/example1.model

	if len(os.Args) < 1+1+2 {
		panic("not enough arguments")
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	numPlayers = n

	// agents
	agents = make([]eval.Agent, 0, len(os.Args[2:]))
	for _, a := range os.Args[2:] {
		agents = append(agents, bot.Parser(a))
	}
}

func main() {

	ds := dataset.LoadDataset("t1k22.json")

	// un tournament reune a varios agentes, y los hace pelear a todos contra todos
	torneo := &eval.Tournament{
		NumPlayers: numPlayers,
		Table:      eval.NewTable(),
		Agents:     agents,
	}

	torneo.Start(ds[:], true)

	torneo.Report()

	// print json
	if bs, err := json.Marshal(torneo.Table); err == nil {
		log.Println(string(bs))
	}

	// guardo el resultado
	t := utils.MiniCurrentTime()
	utils.Write(torneo.Table, "/tmp/res-"+t+".json", true)
	log.Printf("resultado guardado en %s\n\n", "/tmp/res-"+t+".json")

}
