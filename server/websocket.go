package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"net/http"
	"time"
	"fmt"
)

func (s *Server) hubCycle() {
	p := <-s.newPatient
	for {
		time.Sleep(time.Second)
		fmt.Println(*p)
	}
}

//func (s *Server) connectToHub(c *gin.Context) {
//	u := websocket.Upgrader{}
//	conn, err := u.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	// write-only
//	s.hub.list = append(s.hub.list, conn)
//}

func (s *Server) chatWithPatients(c *gin.Context) {
	patientId := c.Param("id")
	if patientId == "" {
		c.String(http.StatusBadRequest, "url must be '/ws_chat/<patient_id>'")
		return
	}

	u := websocket.Upgrader{}
	conn, err := u.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// read messages from assistant/case-manager and send to patient
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.sendResponseToPatient(patientId, string(msg))
		}
	}
}

type HubConn struct {
	list []*websocket.Conn
}

func NewHubConn() *HubConn {
	return &HubConn{
		list: []*websocket.Conn{},
	}
}

type ClientMap struct {
	mtx *sync.RWMutex
	m   map[string]*websocket.Conn
}

func NewClientMap() *ClientMap {
	return &ClientMap{
		mtx: &sync.RWMutex{},
		m:   make(map[string]*websocket.Conn),
	}
}

func (c *ClientMap) set(id string, conn *websocket.Conn) {
	c.mtx.Lock()
	c.m[id] = conn
	c.mtx.Unlock()
}

func (c *ClientMap) get(id string) (*websocket.Conn, bool) {
	c.mtx.RLock()
	conn, ok := c.m[id]
	c.mtx.RUnlock()
	return conn, ok
}
