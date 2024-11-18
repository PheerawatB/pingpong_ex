package main

import (
	"example/pingpong/app/service/player"
	"example/pingpong/app/service/table"
	"log"
)

func main() {

	go func() {
		log.Println("Starting Player Service on port 8888...")
		player.PlayerService() // Start player service on port 8888
	}()

	go func() {
		log.Println("Starting Table Service on port 8889...")
		table.TableService() // Start table service on port 8889
	}()

	// Block the main goroutine to prevent the application from exiting
	select {}

}
