package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var reactURL = "https://slack.com/api/reactions.add"

type SlashResponse struct {
	Challenge string `json:"challenge"`
}

func (s *SlashResponse) respToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

type Slack struct {
	l     *log.Logger
	token string
}

func NewSlack(l *log.Logger, t string) *Slack {
	return &Slack{l, t}
}

type SlackRequest struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		Type        string `json:"type"`
		Channel     string `json:"channel"`
		User        string `json:"user"`
		Text        string `json:"text"`
		Ts          string `json:"ts"`
		EventTs     string `json:"event_ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	Type        string   `json:"type"`
	AuthedTeams []string `json:"authed_teams"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	Challenge   string   `json:"challenge"`
}

type React struct {
	Channel string `json:"channel"`
	Name    string `json:"name"`
	Ts      string `json:"timestamp"`
}

func (s *Slack) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	s.l.Println("Slack endpoint hit")

	if r.Method == http.MethodPost {

		decoder := json.NewDecoder(r.Body)
		var rJSON SlackRequest
		err := decoder.Decode(&rJSON)
		if err != nil {
			panic(err)
		}
		s.l.Println(rJSON.Type)

		if rJSON.Type == "url_verification" {
			rw.Header().Set("Content-Type", "application/json")

			resp := &SlashResponse{
				Challenge: rJSON.Challenge,
			}

			resp.respToJSON(rw)
		}

		if rJSON.Type == "event_callback" {

			if rJSON.Event.Type == "message" {

				if rJSON.Event.User == os.Getenv("SLACKUSER") {

					var bearer = "Bearer " + s.token

					body := &React{
						Channel: rJSON.Event.Channel,
						Name:    "wala",
						Ts:      rJSON.Event.Ts,
					}

					payLoadBuf := new(bytes.Buffer)
					json.NewEncoder(payLoadBuf).Encode(body)
					req, err := http.NewRequest(http.MethodPost, reactURL, payLoadBuf)

					req.Header.Add("Content-Type", "application/json; charset=utf-8")
					req.Header.Add("Authorization", bearer)

					client := &http.Client{}
					res, err := client.Do(req)
					if err != nil {
						s.l.Println("error: ", err)
					}

					defer res.Body.Close()

					respBody, err := io.ReadAll(res.Body)
					if err != nil {
						s.l.Println("Error reading response", err)
					}

					fmt.Println(string(respBody))
				}
			}
		}
	}
}
