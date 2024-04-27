package main

import (
	"keyvaluestore/database"
	"keyvaluestore/routers"
	"log"
)

func main() {
	// Create Gin Instance
	r := routes.SetupRouter()
	//run database
	database.ConnectDB()
	// Run Server
	err := r.Run(":8089")
	if err != nil {
		log.Fatal(err)
	}
}
