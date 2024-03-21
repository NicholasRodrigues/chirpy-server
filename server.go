package main

import (
	"fmt"
	"net/http"
)

func createServer() *http.Server {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()
	apiConfig := apiConfig{
		fileserverHits: 0,
		dbPath:         "chirps.json",
	}

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/reset", apiConfig.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiConfig.validateChirpHandler)
	mux.HandleFunc("POST /api/chirps", apiConfig.insertChirpHandler)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)

	mux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Println("Server created at: " + "http://localhost:" + port + "/app/")
	return server
}
