package post_comment_system

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Comment)
var upgrader = websocket.Upgrader{}

func main() {
	router := mux.NewRouter()

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
	router.Handle("/graphql", h)

	router.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var comment Comment
		err := ws.ReadJSON(&comment)
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcast <- comment
	}
}

func handleMessages() {
	for {
		comment := <-broadcast
		for client := range clients {
			err := client.WriteJSON(comment)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func notifySubscribers(postID string, comment Comment) {
	broadcast <- comment
}
