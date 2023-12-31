package util

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/slack-go/slack"
)

// Message is the struct used to format a slack message
type Message struct {
	message string
	channel string
	status  string
}

func CreateMessage(message string, channel string) Message {
	return Message{
		message: message,
		channel: channel,
		status:  "",
	}
}

func PostMessage(l *log.Logger, m Message) string {
	var (
		retries   int = 3
		errorFlag string
	)

	api := slack.New(os.Getenv("SLACKAPIKEY"))

	// retry slack post until it hits the retry limit or is successful
	for i := 0; i < retries; i++ {
		msgID, _, _, err := api.SendMessage(
			m.channel,
			slack.MsgOptionText(m.message, false),
		)
		if err != nil {
			l.Println("Error posting message: retry:", i, err)
			errorFlag = "fail"
		} else {
			l.Println("Sent message to: ", m.channel, msgID)
			errorFlag = ""
			break
		}

	}

	return errorFlag
}

type SlashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func (s *SlashResponse) RespToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}
