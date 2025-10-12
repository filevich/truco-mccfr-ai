package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/filevich/truco-mccfr-ai/cfr"
	tournament "github.com/filevich/truco-tournament"
	"github.com/filevich/truco-tournament/policies"
	"github.com/filevich/truco-tournament/utils"
	"github.com/truquito/bot/pers"
	"github.com/truquito/bot/schema"
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

func PingN(conn *net.Conn, scanner *bufio.Scanner, n int) {
	// send 3 PING messages
	for i := 0; i < n; i++ {
		utils.Send(conn, "PING")
		msg := ""
		if scanner.Scan() {
			msg = scanner.Text()
		}
		slog.Info("RECEIVED", "i", i, "msg", msg)
		time.Sleep(time.Millisecond * 300)
	}
}

func SetID(conn *net.Conn, scanner *bufio.Scanner, id string) {
	utils.Send(conn, fmt.Sprintf("SET_ID %s", id))
	msg := ""
	if scanner.Scan() {
		msg = scanner.Text()
	}
	slog.Info("RECEIVED", "msg", msg)
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

func StartNewChallenge(conn *net.Conn, scanner *bufio.Scanner) {
	utils.Send(conn, fmt.Sprintf("CHALLENGE %d;2;t1k22;", n))
	msg := ""
	if scanner.Scan() {
		msg = scanner.Text()
	}
	slog.Info("RECEIVED", "msg", msg)
}

func GetName(name string, active bool) string {
	if len(name) == 0 {
		if active {
			name = "activer"
		} else {
			name = utils.RandomString(5)
		}
	}

	// Enforce naming convention: passive agents must start with '_', active agents must not
	if active {
		if strings.HasPrefix(name, "_") {
			panic("Active agents cannot have names starting with '_'")
		}
	} else {
		// Passive agent: ensure '_' prefix
		if !strings.HasPrefix(name, "_") {
			name = "_" + name
		}
	}

	return name
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
	// The scanner is created once and manages its own internal buffer.
	scanner := bufio.NewScanner(*conn)

	defer func() {
		if shutdown {
			utils.Send(conn, "SHUTDOWN")
		} else {
			utils.Send(conn, "EXIT")
		}
	}()

	name = GetName(name, active)
	var (
		p   = &pers.Pers{Nick: name}
		pol = policyFactory(policy)
	)

	// Initial Setup
	// Set ID and wait for the single "OK" response
	utils.Send(conn, fmt.Sprintf("SET_ID %s", name))
	if scanner.Scan() {
		msg := scanner.Text()
		slog.Info("RECEIVED", "msg", msg)
	}

	start := time.Now()

	if active {
		// Start a challenge and wait for the single "OK" response
		utils.Send(conn, fmt.Sprintf("CHALLENGE %d;2;t1k22;", n))
		if scanner.Scan() {
			msg := scanner.Text()
			slog.Info("RECEIVED", "msg", msg)
		}
	}

	// Main Message Processing Loop
	// The scanner.Scan() call will block until a full line (ending in \r\n) is read.
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			slog.Warn("exiting_because_empty_receive")
			return
		}

		slog.Info("RECEIVED", "msg", strings.ReplaceAll(msg, "\"", "'"))

		// NOTE: The original code split by DELIMITER. This is no longer
		// necessary as scanner.Scan() gives us one message at a time.
		// We process the single message `msg`.

		if strings.Contains(msg, "Error") {
			fmt.Println(p.P)
			panic("something weird happened")
		}

		cmd, data := utils.SplitMessage(msg)
		slog.Debug("PROCESSING", "cmd", cmd, "data", strings.ReplaceAll(data, "\"", "'"))

		if cmd == "PONG" || cmd == "OK" {
			// pass
		} else if cmd == "PARSE" || cmd == "MSG" {
			params := strings.Split(data, ";")
			jsonMsgStr := params[0]
			msg := &schema.Msg{}
			err := json.Unmarshal([]byte(jsonMsgStr), msg)
			if err != nil {
				panic(err)
			}
			if err := func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered from panic in Apply: %v", r)
					}
				}()
				p.Apply(msg)
				return nil
			}(); err != nil {
				slog.Error("Failed to apply message", "error", err)
				p.Apply(msg)
				panic(123)
			}
		} else if cmd == "REQ_ACTION" {
			if err := func() (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("recovered from panic in Action: %v", r)
					}
				}()
				mId := strings.TrimSpace(data)
				a := pol.Action(p, mId)
				aStr := a.Stringify()
				txt := fmt.Sprintf("RES_ACTION %s %s", mId, aStr)
				utils.Send(conn, txt)
				return nil
			}(); err != nil {
				slog.Error("Failed to apply message", "error", err)
				utils.Send(conn, "PRINT_STATUS")
				fmt.Println(p.P)
				panic(456)
			}
		} else if cmd == "SERIE_IS_DONE" {
			if !active && !loop {
				return
			}
		} else if cmd == "CHALLENGE_DONE" {
			if loop {
				if active {
					if TimeLimitReached(start) {
						return
					}
					// Start a new challenge
					utils.Send(conn, fmt.Sprintf("CHALLENGE %d;2;t1k22;", n))
					// We expect a single response, so we can call Scan() again here
					if scanner.Scan() {
						response := scanner.Text()
						slog.Info("RECEIVED", "msg", response)
					}
				}
			} else {
				return
			}
		} else {
			fmt.Printf("unknown command: %s\n", cmd)
		}
	}

	// If the loop exits, it's because the connection closed or there was an error.
	if err := scanner.Err(); err != nil {
		slog.Error("Client scanner error", "err", err)
	}
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
