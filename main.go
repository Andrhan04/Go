package main

import (
	"log"
	"net/http"

	"cat-api/config"
	"cat-api/database"
	"cat-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	err := database.Init(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize handlers
	handler := handlers.NewHandler(database.GetDB())

	// Setup router
	router := mux.NewRouter()

	// CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	router.Use(corsMiddleware)

	// Health check
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Cats routes
	catsRouter := router.PathPrefix("/cats").Subrouter()
	catsRouter.HandleFunc("", handler.GetCats).Methods("GET")
	catsRouter.HandleFunc("", handler.CreateCat).Methods("POST")
	catsRouter.HandleFunc("/{id:[0-9]+}", handler.GetCat).Methods("GET")
	catsRouter.HandleFunc("/{id:[0-9]+}", handler.UpdateCat).Methods("PUT")
	catsRouter.HandleFunc("/{id:[0-9]+}", handler.DeleteCat).Methods("DELETE")

	// Types routes
	typesRouter := router.PathPrefix("/types").Subrouter()
	typesRouter.HandleFunc("", handler.GetTypes).Methods("GET")
	typesRouter.HandleFunc("", handler.CreateType).Methods("POST")

	// Masters routes
	mastersRouter := router.PathPrefix("/masters").Subrouter()
	mastersRouter.HandleFunc("", handler.GetMasters).Methods("GET")
	mastersRouter.HandleFunc("", handler.CreateMaster).Methods("POST")
	mastersRouter.HandleFunc("/{id:[0-9]+}", handler.GetMaster).Methods("GET")

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Database file: %s", cfg.DBPath)
	log.Printf("API endpoints:")
	log.Printf("  GET    /health")
	log.Printf("  GET    /cats")
	log.Printf("  POST   /cats")
	log.Printf("  GET    /cats/{id}")
	log.Printf("  PUT    /cats/{id}")
	log.Printf("  DELETE /cats/{id}")
	log.Printf("  GET    /types")
	log.Printf("  POST   /types")
	log.Printf("  GET    /masters")
	log.Printf("  POST   /masters")
	log.Printf("  GET    /masters/{id}")

	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
