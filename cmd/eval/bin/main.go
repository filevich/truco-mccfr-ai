package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/filevich/truco-ai/bot"
	"github.com/filevich/truco-ai/cfr"
	"github.com/filevich/truco-ai/eval"
	"github.com/filevich/truco-ai/eval/dataset"
	"github.com/filevich/truco-ai/utils"
)

var (
	numPlayers             = 0
	agents     []bot.Agent = nil
)

func init() {
	// arg 0 - whole program
	// arg 1 - numPlayers:int
	// args 2..n - model|simple|random

	if len(os.Args) < 1+1+2 {
		panic("not enough arguments")
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	numPlayers = n

	// agents
	agents = make([]bot.Agent, 0, len(os.Args[2:]))
	for i, a := range os.Args[2:] {
		if a == "simple" {
			agents = append(agents, &bot.Simple{})
		} else if a == "random" {
			agents = append(agents, &bot.Random{})
		} else if strings.HasSuffix(a, ".model") {
			name := fmt.Sprintf("m%d", i)
			log.Printf("m%d = %s\n", i, a)
			agent := &cfr.BotCFR{
				N: name,
				F: a,
			}
			agents = append(agents, agent)
		} else {
			panic(fmt.Sprintf("unknown agent `%s`", a))
		}
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
