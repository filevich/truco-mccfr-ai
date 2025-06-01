package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/filevich/truco-mccfr-ai/cfr"
	"github.com/filevich/truco-mccfr-ai/eval"
)

const sep = ":"

func Parser(data string) eval.Agent {
	fields := strings.Split(data, sep)
	agent := strings.ToLower(fields[0])

	if agent == strings.ToLower("simple") {

		return &Simple{}

	} else if agent == strings.ToLower("random") {

		return &Random{}

	} else if agent == strings.ToLower("BotCFR") {

		// BotCFR:<nick>:</path/to/file.model>
		if ok := len(fields) == 3; !ok {
			panic(fmt.Sprintf("bot.Parser: expeted 3 params got %d", len(fields)))
		}
		name, pathToModel := fields[1], fields[2]
		agent := &cfr.BotCFR{
			ID:       name,
			Filepath: pathToModel,
		}
		return agent

	} else if agent == strings.ToLower("BotLazyCFR") {

		// BotCFR:threads:<nick>:</path/to/file.model>
		if ok := len(fields) == 4; !ok {
			panic(fmt.Sprintf("bot.Parser: expeted 4 params got %d", len(fields)))
		}
		threadsStr, name, pathToModel := fields[1], fields[2], fields[3]
		threads, err := strconv.Atoi(threadsStr)
		if err != nil {
			panic(err)
		}
		agent := &cfr.BotLazyCFR{
			ID:       name,
			Filepath: pathToModel,
			Threads:  int64(threads),
		}
		return agent

	} else if agent == strings.ToLower("BotLazyDistilCFR") {

		// BotCFR:threads:<nick>:</path/to/file.model>
		if ok := len(fields) == 4; !ok {
			panic(fmt.Sprintf("bot.Parser: expeted 4 params got %d", len(fields)))
		}
		threadsStr, name, pathToModel := fields[1], fields[2], fields[3]
		threads, err := strconv.Atoi(threadsStr)
		if err != nil {
			panic(err)
		}
		agent := &cfr.BotLazyDistilCFR{
			ID:       name,
			Filepath: pathToModel,
			Threads:  int64(threads),
		}
		return agent

	} else {

		panic(fmt.Sprintf("unknown agent `%s`", data))

	}
}
