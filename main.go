package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/four88/webserver/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
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
	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	mux := http.NewServeMux()
	apiCfg := &apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
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
		id, err := strconv.Atoi(idStg)
		if err != nil {
			responseWithErr(w, "Invalid ID", 404)
		}
		getChirp(w, r, *db, id)
	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
		createUser(w, r, *db)
	})

	mux.Handle("POST /api/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login(w, r, *db, apiCfg.jwtSecret)
	}))

	mux.Handle("PUT /api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleCheckToken(w, r, *db, apiCfg.jwtSecret)
	}))

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
