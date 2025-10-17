package tournamentclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/filevich/truco-mccfr-ai/cfr"
	tournament "github.com/filevich/truco-tournament"
	"github.com/filevich/truco-tournament/policies"
	"github.com/filevich/truco-tournament/utils"
	"github.com/truquito/bot/pers"
	"github.com/truquito/bot/schema"
)

// ChallengeConfig contains configuration for a tournament challenge
type ChallengeConfig struct {
	Active   bool
	Loop     bool
	N        int
	Name     string
	Shutdown bool
	Policy   policies.Policy
}

// Handle manages the tournament client connection and message processing
func Handle(
	conn *net.Conn,
	config ChallengeConfig,
) {
	// The scanner is created once and manages its own internal buffer.
	scanner := bufio.NewScanner(*conn)

	defer func() {
		if config.Shutdown {
			utils.Send(conn, "SHUTDOWN")
		} else {
			utils.Send(conn, "EXIT")
		}
	}()

	name := tournament.GetName(config.Name, config.Active)
	var (
		p   = &pers.Pers{Nick: name}
		pol = config.Policy
	)

	// Initial Setup
	// Set ID and wait for the single "OK" response
	tournament.SetID(conn, scanner, name)

	start := time.Now()

	if config.Active {
		// Start a challenge and wait for the single "OK" response
		tournament.StartNewChallenge(conn, scanner, config.N)
	}

	// Main Message Processing Loop
	// The scanner.Scan() call will block until a full line (ending in \r\n) is read.
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "" {
			slog.Warn("exiting_because_empty_receive")
			return
		}

		slog.Debug("RECEIVED", "msg", strings.ReplaceAll(msg, "\"", "'"))

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
			if !config.Active && !config.Loop {
				return
			}
		} else if cmd == "CHALLENGE_DONE" {
			if config.Loop {
				if config.Active {
					// Check time limit if needed
					// Start a new challenge
					tournament.StartNewChallenge(conn, scanner, config.N)
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
		slog.Error("Client scanner error", "err", err, "start_time", start)
	}
}

// RunChallenge is a simplified wrapper for running a single challenge
// It connects to the tournament server, runs one challenge, and exits
func RunChallenge(addr string, name string, n int, model cfr.ITrainer) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("error connecting to tournament server: %w", err)
	}
	defer conn.Close()

	// Create the CFR bot and policy
	botCFR := &cfr.BotCFR{
		ID:    model.String(),
		Model: model,
	}
	policy := &cfr.CFR_Policy{Model: botCFR}

	config := ChallengeConfig{
		Active:   true,
		Loop:     false,
		N:        n,
		Name:     name,
		Shutdown: false,
		Policy:   policy,
	}

	Handle(&conn, config)
	return nil
}
