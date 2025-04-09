package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Levabu/chirpy/internal/database"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	FileserverHits atomic.Int32
	DB *database.Queries
	Platform string
	JWTSecret string
	PolkaKey string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	
	const filepathRoot = "."
	const port = "8080"

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATORM is not set")
	}	
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	POLKA_KEY := os.Getenv("POLKA_KEY")
	if POLKA_KEY == "" {
		log.Fatal("POLKA_KEY is not set")
	}


	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	cfg := &apiConfig{
		FileserverHits: atomic.Int32{},
		DB: database.New(db),
		Platform: platform,
		JWTSecret: JWTSecret,
		PolkaKey: POLKA_KEY,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	
	mux.HandleFunc("POST /api/login", cfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeToken)

	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhooks)
	
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
	
}
