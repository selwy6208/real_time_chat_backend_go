package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"real-chat-backend/models"
	"real-chat-backend/utils/token"

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
	//current user Id
	userid string
}

// Client management
type ClientManager struct {
	//The client map stores and manages all long connection clients, online is TRUE, and those who are not there are FALSE
	clients map[*Client]bool
	//Web side MESSAGE we use Broadcast to receive, and finally distribute it to all clients
	broadcast   chan []byte
	sendMessage chan []byte
	//Newly created long connection client
	register chan *Client
	//Newly canceled long connection client
	unregister chan *Client
}

// Create a client Manager
var Manager = ClientManager{
	broadcast:   make(chan []byte),
	sendMessage: make(chan []byte),
	register:    make(chan *Client),
	unregister:  make(chan *Client),
	clients:     make(map[*Client]bool),
}

// Will formatting Message into JSON
type SocketMessage struct {
	MessageType string `json:"message_type,omitempty"`
	MessageData string `json:"message_data,omitempty"`
}

type IncomingNewMessage struct {
	MessageType string `json:"message_type,omitempty"`
	MessageData string `json:"message_data,omitempty"`
}
type Message struct {
	Content   string `json:"content,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}

// Will formatting Input Data to get Message into JSON
type InputMessage struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
}
type GetMessagesInput struct {
	ChatUserID uint `json:"chat_user_id" binding:"required"`
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
		case conn := <-Manager.unregister:
			//Determine the state of the connection, if it is true, turn off Send and delete the value of connecting client
			if _, ok := Manager.clients[conn]; ok {
				close(conn.send)
				delete(Manager.clients, conn)
				jsonMessage, _ := json.Marshal(&SocketMessage{MessageType: "close_connection", MessageData: conn.userid})
				Manager.send(jsonMessage, conn)
			}
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
		case message := <-Manager.sendMessage:
			//Traversing the client that has been connected, send the message to them
			var data models.IncomingNewMessage
			err := json.Unmarshal(message, &data)
			if err != nil {
				// Handle the error
				log.Fatal(err.Error())
				return
			}
			for conn := range Manager.clients {
				var _message models.Message
				_err := json.Unmarshal([]byte(data.MessageData), &_message)
				if _err != nil {
					// Handle the error
					log.Fatal(_err.Error())
					return
				}
				if conn.userid == _message.Recipient {
					/* Send the messages to the recipient */
					conn.send <- message
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
		var data models.SocketMessage
		_err := json.Unmarshal(message, &data)
		if _err != nil {
			// Handle the error
			log.Fatal(err.Error())
			return
		}
		if data.MessageType == "new_connection" {
			Manager.send(message, c)
			var incomingNewConnectionUserID string = data.MessageData
			var onlineUsers []string
			for conn := range Manager.clients {
				if conn.id == c.id {
					conn.userid = incomingNewConnectionUserID
				} else {
					if conn.userid != "" {
						onlineUsers = append(onlineUsers, conn.userid)
					}
				}
			}
			/* Get online users */
			resultStr := strings.Join(onlineUsers, ",")
			jsonMessage, _ := json.Marshal(&SocketMessage{MessageType: "online_users", MessageData: resultStr})
			c.send <- jsonMessage

		} else if data.MessageType == "new_message" {
			var message_data string = data.MessageData
			var data models.Message
			_err := json.Unmarshal([]byte(message_data), &data)
			if _err != nil {
				// Handle the error
				log.Fatal(err.Error())
				return
			}
			go SaveMessage(data)
			Manager.sendMessage <- message
		}
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

func SaveMessage(data models.Message) {
	u := models.Message{
		Sender:    data.Sender,
		Recipient: data.Recipient,
		Content:   data.Content,
	}

	_, err := u.SaveMessage()

	if err != nil {
		// Handle the error
		return
	}
}

func GetMessage(c *gin.Context) {
	var input GetMessagesInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myID, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, _ := models.GetMessagesByUserID(myID, input.ChatUserID)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}
