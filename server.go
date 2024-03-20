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

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /healthz", readinessHandler)
	mux.HandleFunc("GET /metrics", apiConfig.metricsHandler)
	mux.HandleFunc("GET /reset", apiConfig.resetHandler)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at port: " + port)
	return server
}
