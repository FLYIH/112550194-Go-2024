package main

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
)

type client chan<- string

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan rxgo.Item) // all incoming client messages
	ObservableMsg = rxgo.FromChannel(messages)
)

func broadcaster() {
	clients := make(map[client]bool)
	MessageBroadcast := ObservableMsg.Observe()
	for {
		select {
		case msg := <-MessageBroadcast:
			for cli := range clients {
				cli <- msg.V.(string)
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn *websocket.Conn, ch <-chan string) {
	for msg := range ch {
		conn.WriteMessage(1, []byte(msg))
	}
}

func wshandle(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "你是 " + who + "\n"
	messages <- rxgo.Of(who + " 來到了現場" + "\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了" + "\n")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		messages <- rxgo.Of(who + " 表示: " + string(msg))
	}
}

var prohibitedWords = map[string]bool{}
var restrictedNames = map[string]bool{}

func readWordsFromFile(filePath string) map[string]bool {
	wordSet := make(map[string]bool)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %s, %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordSet[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s, %v", filePath, err)
	}

	return wordSet
}

func InitObservable() {
	prohibitedWords = readWordsFromFile("swear_word.txt")
	restrictedNames = readWordsFromFile("sensitive_name.txt")
	ObservableMsg = ObservableMsg.
		Filter(func(item interface{}) bool {
			msg, ok := item.(string)
			if !ok {
				log.Println("Filter: Invalid message type, skipping item")
				return false
			}

			for word := range prohibitedWords {
				if strings.Contains(msg, word) {
					log.Println("Filter: Message contains prohibited words, filtering out")
					return false
				}
			}
			return true
		}).
		Map(func(ctx context.Context, item interface{}) (interface{}, error) {
			msg, ok := item.(string)
			if !ok {
				log.Println("Map: Invalid message type, passing through unchanged")
				return item, nil
			}

			for name := range restrictedNames {
				if strings.Contains(msg, name) {
					runes := []rune(name)
					if len(runes) > 1 {
						masked := string(runes[0]) + "*" + string(runes[2:])
						msg = strings.ReplaceAll(msg, name, masked)
					} else {
						msg = strings.ReplaceAll(msg, name, "*")
					}
				}
			}

			sanitized := strings.ToValidUTF8(msg, "")
			if sanitized != msg {
				log.Printf("Map: Sanitized message from %q to %q", msg, sanitized)
			}
			return sanitized, nil
		})
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
