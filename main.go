package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zyr4c31/karaoke/client"
	"github.com/zyr4c31/karaoke/mpv"
	"github.com/zyr4c31/karaoke/server"
)

func main() {
	// init
	if err := godotenv.Load(".env"); err != nil {
		err = errors.New(fmt.Sprintf("godotenv.Load err: %v", err))
		log.Fatal("error loading env file: ", err)
	}
	client.ApiKey = os.Getenv("API_KEY")
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
