package main

import (
	"github.com/flameous/tealtech-go-backend"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	s := tealtech.Server{}
	s.SetDatabase(tealtech.NewDumpDatabase())

	s.Run()
}
