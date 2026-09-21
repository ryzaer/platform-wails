package xhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserCode string
	Conn     *websocket.Conn
	WriteMu  sync.Mutex
}

var clientsMu sync.RWMutex
var clients = map[*websocket.Conn]*Client{}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SocketContext struct {
	Client *Client
	Req    *http.Request
}

func (c *SocketContext) Close() error {
	return c.Client.Conn.Close()
}

func (c *SocketContext) Read(v any) error {

	_, data, err := c.Client.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func write(conn *Client, data []byte) error {

	conn.WriteMu.Lock()
	defer conn.WriteMu.Unlock()

	return conn.Conn.WriteMessage(
		websocket.TextMessage,
		data,
	)
}

func (c *SocketContext) Write(v any) error {

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return write(c.Client, data)
}

func WebSocket(pattern string, handler func(*SocketContext)) {

	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade:", err)
			return
		}

		client := &Client{
			Conn: conn,
		}

		clientsMu.Lock()

		clients[conn] = client
		online := len(clients)

		clientsMu.Unlock()

		log.Println("NEW ONLINE:", online)
		log.Printf("CONNECTED: AS %p VIA %s\n", conn, conn.RemoteAddr())

		handler(&SocketContext{
			Client: client,
			Req:    r,
		})

		clientsMu.Lock()

		client = clients[conn]

		if client != nil {
			log.Println("LEAVE:", client.UserCode)
		}

		delete(clients, conn)

		conn.Close()

		online = len(clients)

		clientsMu.Unlock()

		log.Println("LAST ONLINE:", online)
	})
}

// untuk kirim seluruh user yg online
func Broadcast(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	clientsMu.RLock()

	list := make([]*Client, 0, len(clients))

	for _, client := range clients {
		list = append(list, client)
	}

	clientsMu.RUnlock()

	for _, client := range list {
		write(client, data)
	}
}

func Login(conn *websocket.Conn, code string) {

	// log.Println("############################")
	// log.Println("MASUK KE XHTTP.LOGIN")
	// log.Println("############################")

	clientsMu.Lock()
	defer clientsMu.Unlock()

	if client, ok := clients[conn]; ok {
		client.UserCode = code
	}

	// for c, client := range clients {
	// 	log.Printf("%p -> %s\n", c, client.UserCode)
	// }
}

// untuk kirim user per user
func Send(code string, v any) error {

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	clientsMu.RLock()

	var client *Client

	for _, c := range clients {
		if c.UserCode == code {
			client = c
			break
		}
	}

	clientsMu.RUnlock()

	if client == nil {
		return nil
	}

	return write(client, data)
}
