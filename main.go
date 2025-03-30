package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

var db *sql.DB
var shortcodeRegex = regexp.MustCompile(`^/[\w-]{4,}$`)

func main() {
	var err error

	connStr := "host=localhost port=5432 user=shorty password=secret dbname=shortener sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := setupDatabase(db); err != nil {
		log.Fatalf("Failed to set up DB: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", routeHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodGet:
		w.Write([]byte("👋 Welcome to the URL shortener! Use POST /shorten to shorten a link."))

	case r.URL.Path == "/shorten" && r.Method == http.MethodPost:
		shortenHandler(w, r)

	case shortcodeRegex.MatchString(r.URL.Path):
		redirectHandler(w, r)

	default:
		http.NotFound(w, r)
	}
}

// -- Handler: POST /shorten
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		URL string `json:"url"`
	}
	var body requestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.URL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortcode := generateShortcode(6)

	_, err = db.Exec(
		"INSERT INTO short_links (shortcode, original_url) VALUES ($1, $2)",
		shortcode, body.URL,
	)
	if err != nil {
		http.Error(w, "Could not store URL", http.StatusInternalServerError)
		log.Printf("DB insert error: %v", err)
		return
	}

	type response struct {
		ShortURL string `json:"short_url"`
	}

	resp := response{
		ShortURL: fmt.Sprintf("http://localhost:8080/%s", shortcode),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// -- Handler: GET /{shortcode}
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := r.URL.Path[1:] // Remove the leading slash

	var originalURL string
	err := db.QueryRow("SELECT original_url FROM short_links WHERE shortcode = $1", shortcode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// -- Utility: Setup Table
func setupDatabase(db *sql.DB) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS short_links (
		id SERIAL PRIMARY KEY,
		shortcode TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	return err
}

// -- Utility: Shortcode Generator
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortcode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
