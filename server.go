package main

import (
	"fmt"
	"net/http"
)

func createServer() *http.Server {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()
	apiConfig := apiConfig{}

	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("/reset", apiConfig.resetHandler)
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at port: " + port)
	return server
}
