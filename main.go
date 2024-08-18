package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/four88/webserver/database"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" { // Only count hits for non-readiness paths
			cfg.fileserverHits++
		}
		next.ServeHTTP(w, r) // Call the next handler
	})
}

func main() {
	// Initialize the database
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	mux := http.NewServeMux()
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	// File server handler with metrics middleware
	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fileServer)

	// API endpoints
	mux.HandleFunc("/api/healthz", readinessHandler)
	mux.HandleFunc("/admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		siteHits(w, r, apiCfg)
	})
	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		resetHits(w, r, apiCfg)
	})

	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
			createChirp(w, r, *db)
	})
	
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
			getChirps(w, r, *db)
	})

    mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		idStg := r.PathValue("chirpID")
		id , err := strconv.Atoi(idStg) 
		if(err != nil){
			responseWithErr(w, "Invalid ID", 404)
		}
			getChirp(w, r, *db, id)
	})



	// Start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
