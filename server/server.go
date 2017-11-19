package server

import (
	"github.com/flameous/tealtech-go-backend"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/googollee/go-socket.io"
	"golang.org/x/oauth2/clientcredentials"
	"fmt"
)

type Server struct {
	db         tealtech.Database
	socket     *socketio.Server
	newPatient chan *tealtech.Patient
	botUrl     string
	loginPass  clientcredentials.Config
}

func NewServer() *Server {
	return &Server{
		newPatient: make(chan *tealtech.Patient, 10),
	}
}

func (s *Server) SetDatabase(d tealtech.Database) {
	s.db = d
}

func (s *Server) Run() {
	go s.hubCycle()
	s.socket, _ = socketio.NewServer(nil)

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.HandleMethodNotAllowed = true
	// создание чата с менеджером
	router.POST("/start_new_chat", s.startNewChat)
	// пересыл сообщения пациента из tg-бота на фронт
	router.POST("/send_to_chat", s.sendMessageToManager)

	router.GET("/all_users", s.getAllUsers)
	user := router.Group("/user")
	user.GET("/:id", s.getUser)
	user.POST("/:id", s.saveUser)

	router.GET("/ws_hub/", s.socketHandler)
	router.POST("/ws_hub/", s.socketHandler)
	router.Handle("WS", "/ws_hub/", s.socketHandler)
	router.Handle("WSS", "/ws_hub/", s.socketHandler)

	log.Fatal(router.Run(":8100"))
}

func (s *Server) socketHandler(c *gin.Context) {
	s.socket.On("connection", func(so socketio.Socket) {

		so.On("say", func(msg string) {
			fmt.Println(msg)
			// fmt.Println("emit:", so.Emit("chat message", msg))
			// so.BroadcastTo("say", "chat message", msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	s.socket.On("error", func(so socketio.Socket, err error) {
		log.Printf("[ WebSocket ] Error : %v", err.Error())
	})

	s.socket.ServeHTTP(c.Writer, c.Request)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*.ngrok.io/*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
