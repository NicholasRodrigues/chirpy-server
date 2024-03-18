package main

import (
	"fmt"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func createServer() *http.Server {
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/app", http.FileServer(http.Dir("./app/"))))
	mux.HandleFunc("/healthz", readinessHandler)
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at port: " + port)
	return server
}

func main() {
	server := createServer()
	err := server.ListenAndServe()
	fmt.Println("Server started")
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
