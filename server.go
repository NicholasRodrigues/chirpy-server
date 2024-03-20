package main

import (
	"fmt"
	"net/http"
)

func createServer() *http.Server {
	const port = "8080"
	mux := http.NewServeMux()
	apiConfig := apiConfig{}

	mux.Handle("/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app/")))))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.Handle("/metrics", http.HandlerFunc(apiConfig.metricsHandler))
	mux.Handle("/reset", http.HandlerFunc(apiConfig.resetHandler))
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at port: " + port)
	return server
}
