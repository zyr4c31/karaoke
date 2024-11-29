package main

import (
	"log"

	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
	"github.com/zyr4c31/karaoke/server"
)

func main() {
	// init
	if err := client.GetApiKey(); err != nil {
		log.Panic(err)
	}
	// init

	conn, err := mpv.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	if err := server.Run(conn); err != nil {
		log.Panic(err)
	}
}
