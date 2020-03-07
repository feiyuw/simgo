package server

import (
	"time"
)

const (
	MSGSIZE = 1000
)

type Message struct {
	Method    string `json:"method"`
	Direction string `json:"direction"`
	From      string `json:"from"`
	To        string `json:"to"`
	Ts        int64  `json:"ts"`
	Body      string `json:"body"`
}

func newMessageRecorder(server *Server) func(mtd, direction, from, to, body string) error {
	server.Messages = make([]*Message, 0, MSGSIZE)
	return func(mtd, direction, from, to, body string) error {
		msg := &Message{
			Method:    mtd,
			Direction: direction,
			From:      from,
			To:        to,
			Ts:        time.Now().UnixNano() / int64(time.Millisecond),
			Body:      body,
		}
		// TODO: add rlock
		if len(server.Messages) == MSGSIZE {
			copy(server.Messages[1:], server.Messages[0:MSGSIZE-1])
			server.Messages[0] = msg
		} else {
			server.Messages = append([]*Message{msg}, server.Messages...)
		}
		return nil
	}
}

func queryMessages(server *Server, skip, limit int) []*Message {
	msgCnt := len(server.Messages)

	if msgCnt <= skip {
		return []*Message{}
	} else if msgCnt <= skip+limit {
		return server.Messages[skip:msgCnt]
	} else {
		return server.Messages[skip : skip+limit]
	}
}
