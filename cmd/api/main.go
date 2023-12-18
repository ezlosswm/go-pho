package main

import (
	"log"
	"os"
	"pb/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.NewServer()

	port := os.Getenv("PORT")
	log.Println("listening on port: ", port)
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
