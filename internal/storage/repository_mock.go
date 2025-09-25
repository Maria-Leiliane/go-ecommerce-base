package storage

import "e-commerce.com/internal/domain"

// MockProductRepository is a mock implementation of the ProductRepository interface.
type MockProductRepository struct {
	Products []domain.Product
	Error    error
}

func (m *MockProductRepository) Save(product *domain.Product) error {
	if m.Error != nil {
		return m.Error
	}
	product.ID = len(m.Products) + 1
	m.Products = append(m.Products, *product)
	return nil
}

func (m *MockProductRepository) FindAll() ([]domain.Product, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Products, nil
}
func (m *MockProductRepository) FindByID(_ int) (domain.Product, error) {
	return domain.Product{}, m.Error
}
func (m *MockProductRepository) Update(_ *domain.Product) error { return m.Error }
func (m *MockProductRepository) Delete(_ int) error             { return m.Error }
