// ecommerce-base/cmd/api/e2e_test.go
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"e-commerce.com/internal/domain"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB
var testServer *httptest.Server

// Special function that manages the lifecycle of tests in this package.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start the PostgreSQL container for testing
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("could not start postgres container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("could not terminate postgres container: %s", err)
		}
	}()

	// Get the connection string and connect to the test database
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("could not get connection string: %s", err)
	}
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to test database: %s", err)
	}

	// Run the migrations
	createTableSQL := `CREATE TABLE products (
		id SERIAL PRIMARY KEY, name TEXT NOT NULL, price NUMERIC(10, 2) NOT NULL,
		amount INTEGER NOT NULL, description TEXT
	);`
	if _, err = testDB.Exec(createTableSQL); err != nil {
		log.Fatalf("could not create table: %s", err)
	}

	// Start a test server using the router on a random port.
	router := setupRouter(testDB)
	testServer = httptest.NewServer(router)
	defer testServer.Close()

	exitCode := m.Run()

	os.Exit(exitCode)
}

// Testing the complete life cycle of a product via API.
func TestE2E_ProductLifecycle(t *testing.T) {
	var createdProduct domain.Product

	// --- 1. CREATE (POST /products) ---
	t.Run("Create a new product", func(t *testing.T) {
		productJSON := `{"name": "E2E Test Mouse", "price": 150.75, "amount": 20, "description": "A mouse for E2E testing"}`
		body := bytes.NewBufferString(productJSON)

		resp, err := http.Post(testServer.URL+"/products", "application/json", body)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("Failed to Close: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201 Created, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdProduct); err != nil {
			t.Fatalf("Failed to decode created product: %v", err)
		}

		if createdProduct.ID == 0 {
			t.Error("Expected product ID to be set, but it was 0")
		}
		if createdProduct.Name != "E2E Test Mouse" {
			t.Errorf("Expected product name to be 'E2E Test Mouse', got '%s'", createdProduct.Name)
		}
	})

	// READ (GET /products/{id})
	t.Run("Get the created product", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/products/%d", testServer.URL, createdProduct.ID))
		if err != nil {
			t.Fatalf("Failed to get product: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("Failed to Close: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var fetchedProduct domain.Product
		if err := json.NewDecoder(resp.Body).Decode(&fetchedProduct); err != nil {
			t.Fatalf("Failed to decode fetched product: %v", err)
		}

		if fetchedProduct.ID != createdProduct.ID || fetchedProduct.Name != createdProduct.Name {
			t.Errorf("Fetched product does not match created product. Got %+v", fetchedProduct)
		}
	})

	// UPDATE (PUT /products/{id})
	t.Run("Update the product", func(t *testing.T) {
		updatedJSON := `{"name": "Updated E2E Mouse", "price": 160.00, "amount": 15, "description": "An updated mouse"}`
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/products/%d", testServer.URL, createdProduct.ID), bytes.NewBufferString(updatedJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to update product: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatalf("Failed to Close: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
		}

		var updatedProduct domain.Product
		if err := json.NewDecoder(resp.Body).Decode(&updatedProduct); err != nil {
			t.Fatalf("Failed to decode updated product: %v", err)
		}

		if updatedProduct.Name != "Updated E2E Mouse" || updatedProduct.Amount != 15 {
			t.Errorf("Product was not updated correctly. Got %+v", updatedProduct)
		}
	})

	// DELETE (DELETE /products/{id})
	t.Run("Delete the product", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/products/%d", testServer.URL, createdProduct.ID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete product: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("Failed to Close: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK { // Or 204 if the handler changes
			t.Fatalf("Expected status 200 OK, got %d", resp.StatusCode)
		}
	})

	// VERIFY DELETION (GET /products/{id})
	t.Run("Verify product is deleted", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/products/%d", testServer.URL, createdProduct.ID))
		if err != nil {
			t.Fatalf("Failed to get product after deletion: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatalf("Failed to Close: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status 404 Not Found after deletion, got %d", resp.StatusCode)
		}
	})
}
