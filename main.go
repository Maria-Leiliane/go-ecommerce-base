// go-ecommerce-base/main.go
package main

// @title           Products API
// @version         1.0
// @description     API CRUD of Products with Go.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	productHandler "e-commerce.com/internal/handler/http"
	"e-commerce.com/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// setupRouter creates and configures the chi router with all dependencies and routes.
func setupRouter(db *sql.DB) *chi.Mux {
	productRepo := storage.NewProductRepository(db)
	productH := productHandler.NewProductHandler(productRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productH.ListProducts)
		r.Post("/", productH.CreateProduct)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productH.GetProduct)
			r.Put("/", productH.UpdateProduct)
			r.Delete("/", productH.DeleteProduct)
		})
	})

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not load .env file.")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Warning: Could not Close the Database.")
		}
	}(db)

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY, name TEXT NOT NULL, price NUMERIC(10, 2) NOT NULL,
		amount INTEGER NOT NULL, description TEXT
	);`
	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Database connected and table ready.")

	// Just call setupRouter and start the server.
	router := setupRouter(db)
	log.Println("Server starting on port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
