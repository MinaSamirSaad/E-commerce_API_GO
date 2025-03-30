package products

import (
	"database/sql"

	"github.com/MinaSamirSaad/ecommerce/services/shared"
)

type ProductStore interface {
	GetProducts() ([]shared.Product, error)
	CreateProduct(p *shared.CreateProductPayload) (*shared.Product, error)
	GetProductByID(id int) (*shared.Product, error)
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]shared.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := []shared.Product{}
	for rows.Next() {
		p := shared.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *Store) CreateProduct(p *shared.CreateProductPayload) (*shared.Product, error) {
	result, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", p.Name, p.Description, p.Image, p.Price, p.Quantity)
	if err != nil {
		return nil, err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the product by ID
	return s.GetProductByID(int(lastID))
}

func (s *Store) GetProductByID(id int) (*shared.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	p := &shared.Product{}
	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Image, &p.Price, &p.Quantity, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	if p.ID == 0 {
		return nil, nil
	}
	return p, nil
}
