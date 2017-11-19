package server

import (
	"github.com/flameous/tealtech-go-backend"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"net/url"
	"log"
)

func (s *Server) getUser(c *gin.Context) {
	u := s.db.GetPatient(c.Param("id"))
	if u == nil {
		c.String(http.StatusNotFound, "not found")
	} else {
		c.JSON(http.StatusOK, u)
	}
}

func (s *Server) getAllUsers(c *gin.Context) {
	c.JSON(200, s.db.GetAllPatients())
}

func (s *Server) saveUser(c *gin.Context) {
	userData, ok := c.GetPostForm("patient")
	if !ok {
		c.String(http.StatusBadRequest, "missing 'patient' field!")
		return
	}

	fmt.Println(userData)
	p := new(tealtech.Patient)
	if err := json.Unmarshal([]byte(userData), p); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	s.db.SavePatient(c.Param("id"), p)
	c.String(http.StatusOK, `user saved`)
}

func (s *Server) startNewChat(c *gin.Context) {
	uid, ok := c.GetPostForm("uid")
	if !ok {
		c.String(http.StatusBadRequest, `missing "uid"`)
		return
	}

	u := s.db.GetPatient(uid)
	if u != nil {
		s.newPatient <- u
		c.String(http.StatusOK, "ok, wait pls")
	} else {
		c.String(http.StatusNotFound, "chat not found")
	}
}

// message from patient to manager/assistant
func (s *Server) sendMessageToManager(c *gin.Context) {
	//uid, ok := c.GetPostForm("uid")
	//if !ok {
	//	c.String(http.StatusBadRequest, "missing 'uid'")
	//	return
	//}
	//text, ok := c.GetPostForm("text")
	//if !ok {
	//	c.String(http.StatusBadRequest, `missing "text"`)
	//	return
	//}
	//conn, ok := s.chats.get(uid)
	//if !ok {
	//	c.String(http.StatusNotFound, "not found")
	//	return
	//}
	//
	//conn.WriteMessage(websocket.TextMessage, []byte(text))
	//c.String(http.StatusOK, "ok")
}

func (s *Server) sendResponseToPatient(id, message string) bool {
	vals := url.Values{}
	vals.Set("uid", id)
	vals.Set("message", message)
	resp, err := http.PostForm(s.botUrl, vals)
	if err != nil {
		log.Println(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	} else {
		return false
	}
}

func (s *Server) reset(c *gin.Context) {
	s.db.Reset()
	c.String(http.StatusOK, "reset")
}

func (s *Server) anyJira(c *gin.Context) {
	fmt.Println(c.Param("id"))
	b, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("Jira Response:\n" + string(b))
	c.String(200, "ok")
}


