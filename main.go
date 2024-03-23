package main

import (
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err2 := godotenv.Load()
	if err2 != nil {
		return
	}
	server := createServer()
	err := server.ListenAndServe()
	fmt.Println("Server started")
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
