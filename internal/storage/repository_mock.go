package storage

import (
	"fmt"

	"e-commerce.com/internal/domain"
)

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

func (m *MockProductRepository) FindAll(page, limit int) ([]domain.Product, int, error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	total := len(m.Products)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		return []domain.Product{}, total, nil
	}
	if end > total {
		end = total
	}
	return m.Products[start:end], total, nil
}

func (m *MockProductRepository) FindByID(id int) (domain.Product, error) {
	if m.Error != nil {
		return domain.Product{}, m.Error
	}
	for _, p := range m.Products {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Product{}, fmt.Errorf("product not found")
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	if m.Error != nil {
		return m.Error
	}
	for i, p := range m.Products {
		if p.ID == product.ID {
			m.Products[i] = *product
			return nil
		}
	}
	return fmt.Errorf("product not found for update")
}

func (m *MockProductRepository) Delete(id int) error {
	if m.Error != nil {
		return m.Error
	}
	idx := -1
	for i, p := range m.Products {
		if p.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("product not found for deletion")
	}
	m.Products = append(m.Products[:idx], m.Products[idx+1:]...)
	return nil
}
