package tealtech

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"fmt"
	"io/ioutil"
)

type Server struct {
	db Database
}

func (s *Server) SetDatabase(d Database) {
	s.db = d
}

func (s *Server) getUser(c *gin.Context) {
	login := c.Query("login")
	uid := c.Query("uid")
	id, _ := strconv.Atoi(uid)

	u := s.db.GetUser(id, login)
	if u == nil {
		c.String(http.StatusNotFound, "not found")
	} else {
		c.JSON(http.StatusOK, u)
	}
}

func (s *Server) reset(c *gin.Context) {
	s.db.Reset()
	c.String(http.StatusOK, "reset")
}

func (s *Server) AnyShit(c *gin.Context) {
	b, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(b))
	c.String(200, "ok")
}

func (s *Server) saveUser(c *gin.Context) {
	u := new(BotUser)
	raw, ok := c.GetPostForm("user")
	if !ok {
		c.String(http.StatusBadRequest, "missing 'user' field!")
		return
	}

	if err := json.Unmarshal([]byte(raw), u); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	s.db.SaveUser(u)
	c.String(200, `user saved`)
}

func (s *Server) Run() {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	bot := router.Group("/bot")
	bot.GET("/get_user", s.getUser)
	bot.POST("/save_user", s.saveUser)
	bot.POST("/reset", s.reset)

	jira := router.Group("/jira")
	jira.Any("/*any", s.AnyShit)
	log.Fatal(router.Run(":8100"))
}
