package domain

import _ "database/sql"

// Product defines the structure for a product item.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Amount      int     `json:"amount"`
	Description string  `json:"description"`
}

type ProductRepository interface {
	Save(product *Product) error
	FindAll(page, limit int) ([]Product, int, error)
	FindByID(id int) (Product, error)
	Update(product *Product) error
	Delete(id int) error
}
