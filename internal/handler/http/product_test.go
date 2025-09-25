package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"e-commerce.com/internal/domain"
	"e-commerce.com/internal/storage"
)

func TestListProductsHandler(t *testing.T) {
	// 1. Set up
	// The mock with predefined data.
	mockRepo := &storage.MockProductRepository{
		Products: []domain.Product{
			{ID: 1, Name: "Test Product", Price: 10.0, Amount: 5, Description: "A test product"},
		},
	}
	productHandler := NewProductHandler(mockRepo)

	// The request HTTP test.
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	// "Response Recorder" to record the handler's response.
	rr := httptest.NewRecorder()

	// 2. Execution
	// Call the handler directly.
	productHandler.ListProducts(rr, req)

	// 3. Assertions
	// Check the status code 200 OK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the expected data is in the response body.
	expected := `[{"id":1,"name":"Test Product","price":10,"amount":5,"description":"A test product"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
