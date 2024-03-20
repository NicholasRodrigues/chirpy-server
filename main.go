package main

import (
	"fmt"
)

func main() {
	server := createServer()
	err := server.ListenAndServe()
	fmt.Println("Server started")
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
