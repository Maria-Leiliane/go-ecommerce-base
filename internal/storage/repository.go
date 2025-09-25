package storage

import (
	"database/sql"
	"errors"
	"log"

	"e-commerce.com/internal/domain"
)

// pgProductRepository implements the ProductRepository interface for PostgreSQL.
type pgProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new instance of the product repository.
func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &pgProductRepository{db: db}
}

func (r *pgProductRepository) Save(product *domain.Product) error {
	sqlStatement := `INSERT INTO products (name, price, amount, description) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(sqlStatement, product.Name, product.Price, product.Amount, product.Description).Scan(&product.ID)
}

// FindAll now accepts page and limit, and returns the product slice, total count, and an error.
func (r *pgProductRepository) FindAll(page, limit int) ([]domain.Product, int, error) {
	var total int
	// First, get the total count of products.
	err := r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate the offset for pagination.
	offset := (page - 1) * limit

	// Now, fetch the products for the specific page.
	rows, err := r.db.Query("SELECT id, name, price, amount, description FROM products ORDER BY id ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows on FindAll: %v", err)
		}
	}(rows)

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Amount, &p.Description); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *pgProductRepository) FindByID(id int) (domain.Product, error) {
	row := r.db.QueryRow("SELECT id, name, price, amount, description FROM products WHERE id = $1", id)
	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Amount, &p.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Product{}, errors.New("product not found")
		}
		return domain.Product{}, err
	}
	return p, nil
}

func (r *pgProductRepository) Update(product *domain.Product) error {
	sqlStatement := `UPDATE products SET name=$1, price=$2, amount=$3, description=$4 WHERE id=$5`
	res, err := r.db.Exec(sqlStatement, product.Name, product.Price, product.Amount, product.Description, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found for update")
	}
	return nil
}

func (r *pgProductRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found for deletion")
	}
	return nil
}
