package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	MessageAddContainer = iota + 1
	MessageRemoveContainer
	MessageFetchContainers
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Socket struct {
	conn              *websocket.Conn
	outgres           chan []byte
	OnClientConnected func()
}

func newSocket() *Socket {
	return &Socket{
		conn:    nil,
		outgres: make(chan []byte, 10),
	}
}

func (s *Socket) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	var err error
	s.conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}

	if s.OnClientConnected != nil {
		s.OnClientConnected()
	}

	go s.readFromClient()
	go s.writeToClient()
}

func (s *Socket) Send(message interface{}) {
	data, _ := json.Marshal(message)
	s.outgres <- data
}

func (s *Socket) readFromClient() {
	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			break
		}
		log.Println(string(data))
	}
}

func (s *Socket) writeToClient() {
	for data := range s.outgres {
		s.conn.WriteMessage(websocket.TextMessage, data)
	}
	s.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
