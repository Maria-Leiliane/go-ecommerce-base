package http

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"

	"e-commerce.com/internal/domain"

	"github.com/go-chi/chi/v5"
)

// ProductHandler Definition of the ProductHandler struct.
type ProductHandler struct {
	repo domain.ProductRepository
}

// PaginatedResponse is the structure for the paginated response.
type PaginatedResponse struct {
	Data        []domain.Product `json:"data"`
	TotalPages  int              `json:"total_pages"`
	CurrentPage int              `json:"current_page"`
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(repo domain.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Aux functions

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
// @Summary      List all products with pagination
// @Description  Returns a paginated list of all products.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(50)
// @Success      200    {object}  PaginatedResponse
// @Failure      500    {object}  map[string]string
// @Router       /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 50 {
		limit = 50
	}

	products, total, err := h.repo.FindAll(page, limit)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		log.Printf("Error finding all products: %v", err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := PaginatedResponse{
		Data:        products,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	h.respondWithJSON(w, http.StatusOK, response)
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
