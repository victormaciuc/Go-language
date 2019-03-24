package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func main(){
	serveFrontend()
	serveBackend()

	err := http.ListenAndServe("localhost:1324", nil)
	if err != nil{
		log.Fatalf("could not start webserver: %v", err)
	}
}

	func serveFrontend(){
		fs := http.FileServer(http.Dir("./messenger-ui/build/"))
		http.Handle("/", fs)
	}

	var clients []*Client

	type Client struct {
		//ID int
		Conn *websocket.Conn

	}

	func serveBackend(){
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
			upgrader := websocket.Upgrader{}

			client := &Client{}

			conn, err := upgrader.Upgrade(w,r,nil)
			if err != nil{
				log.Fatalf("could not upgrade: %v", err)
			}

			client.Conn = conn
			clients = append(clients, client)

			err = client.Conn.WriteMessage(websocket.TextMessage, []byte("{\"Text\":\"First message\"}"))
			err = client.Conn.WriteMessage(websocket.TextMessage, []byte("{\"Text\":\"Second message\"}"))

			if err != nil{
				log.Fatalf("could not send message: %v", err)
			}

			var wg sync.WaitGroup

			wg.Add(5)

			go func(wg sync.WaitGroup) {

				//asculta dupa mesaje thread1
				for {
					messageType, msg, _ := client.Conn.ReadMessage()
					fmt.Printf("message received: %s (%d message type)", msg, messageType)

					for _, client := range clients {
						_ = client.Conn.WriteMessage(messageType, msg)
					}
					if false == true{
						wg.Done()
						break
					}
				}
			}(wg)

			wg.Add(5)
			//sms thread2
			go func() {
				for {
					fmt.Printf("Working...")
					wg.Done()
				}
			}()
			wg.Wait()

		})
	}