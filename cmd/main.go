package main

import (
	"github.com/flameous/tealtech-go-backend"
	"github.com/flameous/tealtech-go-backend/server"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	s := server.NewServer()
	//s.SetDatabase(tealtech.NewDumpDatabase())
	s.SetDatabase(tealtech.NewAidBoxDatabase())
	s.Run()
}
