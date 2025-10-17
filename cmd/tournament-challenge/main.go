package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/filevich/truco-mccfr-ai/cfr"
	"github.com/filevich/truco-mccfr-ai/internal/tournamentclient"
	tournament "github.com/filevich/truco-tournament"
	"github.com/filevich/truco-tournament/policies"
	"github.com/filevich/truco-tournament/utils"
)

var (
	tournamentAddr = flag.String("tournament_addr", "localhost:8080", "address of the tournament")
	active         = flag.Bool("active", true, "run in active mode")
	loop           = flag.Bool("loop", false, "keep connectd even if the serie is done")
	policy         = flag.String("policy", "random", "Bot policy")
	n              = flag.Int("n", 1000, "Number of (doubled) games")
	name           = flag.String("name", "", "Bot name")
	level          = flag.String("level", "warn", "set the logging level (debug, info, warn, error)")
	shutdown       = flag.Bool("shutdown", false, "send SHUTDOWN message (instead of an EXIT one) at the end")
)

func init() {
	flag.Parse()
	tournament.LogSystemsVersion()
	// initialize logging
	slogLevel, err := utils.ParseLogLevel(*level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level: %s\n", err)
		os.Exit(1)
	}
	utils.InitLogs(slogLevel)
}

func policyFactory(policy string) policies.Policy {
	var pol policies.Policy

	modelPath := os.Getenv("model")
	if modelPath == "" {
		panic("environment variable 'model' is required for CFR policy")
	}
	model := &cfr.BotCFR{
		ID:       "?",
		Filepath: modelPath,
	}
	model.Initialize()

	switch policy {
	case "cfr":
		pol = &cfr.CFR_Policy{Model: model}
	case "cfr_greedy":
		meta_model := &cfr.BotCFR_Greedy{BotCFR: *model}
		pol = &cfr.CFR_Policy_Greedy{Model: meta_model}
	default:
		panic(fmt.Sprintf("Unknown policy %s", policy))
	}

	return pol
}

func TimeLimitReached(start time.Time) bool {
	limit_sec_str := os.Getenv("LIMIT")
	if len(limit_sec_str) == 0 {
		limit_sec_str = "300"
	}

	limit_sec, err := strconv.Atoi(limit_sec_str)
	runtime := time.Since(start).Seconds()
	slog.Info("PROGRESS", "limit", float64(limit_sec), "runtime", runtime)
	if err == nil && runtime > float64(limit_sec) {
		slog.Warn("TIME_LIMIT_REACHED")
		return true
	}

	return false
}

func handle(
	conn *net.Conn,
	active bool,
	loop bool,
	policy string,
	n int,
	name string,
	shutdown bool,
) {
	pol := policyFactory(policy)

	config := tournamentclient.ChallengeConfig{
		Active:   active,
		Loop:     loop,
		N:        n,
		Name:     name,
		Shutdown: shutdown,
		Policy:   pol,
	}

	tournamentclient.Handle(conn, config)
}

func main() {
	conn, err := net.Dial("tcp", *tournamentAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	handle(&conn, *active, *loop, *policy, *n, *name, *shutdown)
	time.Sleep(time.Second)
}
