package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client
type Client struct {
	//User ID
	id string
	//Connected socket
	socket *websocket.Conn
	//Message
	send chan []byte
}

// Client management
type ClientManager struct {
	//The client map stores and manages all long connection clients, online is TRUE, and those who are not there are FALSE
	clients map[*Client]bool
	//Web side MESSAGE we use Broadcast to receive, and finally distribute it to all clients
	broadcast chan []byte
	//Newly created long connection client
	register chan *Client
	//Newly canceled long connection client
	unregister chan *Client
}

// Create a client Manager
var Manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

// Will formatting Message into JSON
type Message struct {
	//Message Struct
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (Manager *ClientManager) Start() {
	for {
		select {
		//If there is a new connection access, pass the connection to conn through the channel
		case conn := <-Manager.register:
			//Set the client connection to true
			Manager.clients[conn] = true
			//Format the message of returning to the successful connection JSON
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected. "})
			//Call the client's send method and send messages
			Manager.send(jsonMessage, conn)
			//If the connection is disconnected
		case conn := <-Manager.unregister:
			//Determine the state of the connection, if it is true, turn off Send and delete the value of connecting client
			if _, ok := Manager.clients[conn]; ok {
				close(conn.send)
				delete(Manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected. "})
				Manager.send(jsonMessage, conn)
			}
			//broadcast
		case message := <-Manager.broadcast:
			//Traversing the client that has been connected, send the message to them
			for conn := range Manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(Manager.clients, conn)
				}
			}
		}
	}
}

// Define the send method of client management
func (Manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range Manager.clients {
		//Send messages not to the shielded connection
		if conn != ignore {
			conn.send <- message
		}
	}
}

// Define the read method of the client structure
func (c *Client) read() {
	defer func() {
		Manager.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		//Read message
		_, message, err := c.socket.ReadMessage()
		//If there is an error message, cancel this connection and then close it
		if err != nil {
			Manager.unregister <- c
			_ = c.socket.Close()
			break
		}
		//If there is no error message, put the information in Broadcast
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		Manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for {
		select {
		//Read the message from send
		case message, ok := <-c.send:
			//If there is no message
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//Write it if there is news and send it to the web side
			_ = c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func WsHandler(c *gin.Context) {
	//Upgrade the HTTP protocol to the websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//Every connection will open a new client, client.id generates through UUID to ensure that each time it is different
	client := &Client{id: uuid.Must(uuid.NewV4(), nil).String(), socket: conn, send: make(chan []byte)}
	//Register a new link
	Manager.register <- client

	//Start the message to collect the news from the web side
	go client.read()
	//Start the corporation to return the message to the web side
	go client.write()
}
