package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"e-commerce.com/internal/domain"

	"github.com/go-chi/chi/v5"
)

// ProductHandler holds the dependencies for the product handlers.
type ProductHandler struct {
	repo domain.ProductRepository
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(repo domain.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// --- Helper Functions ---

func (h *ProductHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *ProductHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// --- Handler Methods ---

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product based on the provided JSON payload. The created product, including its new ID, is returned.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      domain.Product  true  "Product Payload"
// @Success      201      {object}  domain.Product
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if p.Name == "" || p.Price < 0 || p.Amount < 0 {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product data: name, price, and amount are required and must be valid")
		return
	}

	if err := h.repo.Save(&p); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		log.Printf("Error saving product: %v", err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, p)
}

// ListProducts godoc
// @Summary      List all products
// @Description  Returns a list of all products registered in the database.
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      500  {object}  map[string]string
// @Router       /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, _ *http.Request) {
	products, err := h.repo.FindAll()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		log.Printf("Error finding all products: %v", err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, products)
}

// GetProduct godoc
// @Summary      Get a product by ID
// @Description  Retrieves the details of a specific product by its unique ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.repo.FindByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			h.respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve product")
			log.Printf("Error finding product by ID: %v", err)
		}
		return
	}
	h.respondWithJSON(w, http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Updates the details of an existing product identified by its ID using the provided JSON payload.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Product ID"
// @Param        product  body      domain.Product  true  "Product Payload"
// @Success      200      {object}  domain.Product
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	p.ID = id
	if err := h.repo.Update(&p); err != nil {
		if err.Error() == "product not found for update" {
			h.respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to update product")
			log.Printf("Error updating product: %v", err)
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, p)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes a product from the database by its unique ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "product not found for deletion" {
			h.respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
			log.Printf("Error deleting product: %v", err)
		}
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
