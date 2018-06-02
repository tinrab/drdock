package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	socket := newSocket()

	app, err := newApp(func(action Action, data interface{}) {
		switch action {
		case ActionContainerStarted:
			socket.Send(newMessage(MessageAddContainer, data))
		case ActionContainerStopped:
			socket.Send(newMessage(MessageRemoveContainer, data))
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	socket.OnClientConnected = func() {
		log.Println("Client connected")
		containers := []Container{}
		for _, c := range app.Containers {
			containers = append(containers, c)
		}
		socket.Send(newMessage(MessageFetchContainers, containers))
	}

	go app.Run(ctx)
	http.HandleFunc("/ws", socket.HandleWebSocket)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func newMessage(kind int, payload interface{}) interface{} {
	return struct {
		Kind    int         `json:"kind"`
		Payload interface{} `json:"payload"`
	}{
		kind,
		payload,
	}
}
