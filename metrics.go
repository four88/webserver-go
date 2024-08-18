package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func htmlTemplate(hit int) string {
	htmlTemplate := fmt.Sprintf(`
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, hit)
	return htmlTemplate
}


func siteHits(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(
		http.StatusOK,
	)
	html := htmlTemplate(cfg.fileserverHits)
	w.Write([]byte(html))
  // w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileserverHits)))
}

func resetHits(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(
		http.StatusOK,
	)
	cfg.fileserverHits = 0
	w.Write([]byte("Hits: " + strconv.Itoa(cfg.fileserverHits)))
}
