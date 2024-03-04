package handlers

import (
	"net/http"
	"time"
)

var Clients = make([]chan string, 0)

func HandleGetSubscribe(w http.ResponseWriter, r *http.Request) {
	channel := make(chan string)
	pos := len(Clients)
	Clients = append(Clients, channel)

	select {
	case <-time.After(30 * time.Second):
	case message := <-channel:
		if len(Clients) > 1 {
			Clients = append(Clients[:pos], Clients[pos+1:]...)
		} else {
			Clients = make([]chan string, 0)
		}
		w.Write([]byte(message))
	}
}
