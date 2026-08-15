package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"ecommerce-backend/internal/models"
)

type Store struct {
	mu            sync.RWMutex
	products      []models.Product
	orders        []models.Order
	nextProductID int
	dataPath      string
}
type data struct {
	Products      []models.Product `json:"products"`
	Orders        []models.Order   `json:"orders"`
	NextProductID int              `json:"nextProductId"`
}

func New(dataPath string) *Store {
	s := &Store{dataPath: dataPath, nextProductID: 4, products: []models.Product{
		{ID: 1, Name: "Everyday Tote", Description: "A clean, durable carryall for daily essentials.", Category: "Bags", Price: 79, Stock: 24, Active: true},
		{ID: 2, Name: "Minimal Watch", Description: "A refined timepiece with a comfortable leather strap.", Category: "Accessories", Price: 129, Stock: 12, Active: true},
		{ID: 3, Name: "Studio Lamp", Description: "Warm directional lighting for a focused workspace.", Category: "Home", Price: 95, Stock: 8, Active: true},
	}}
	s.load()
	return s
}
func (s *Store) load() {
	file, err := os.Open(s.dataPath)
	if err != nil {
		return
	}
	defer file.Close()
	var saved data
	if json.NewDecoder(file).Decode(&saved) == nil && saved.NextProductID > 0 {
		s.products, s.orders, s.nextProductID = saved.Products, saved.Orders, saved.NextProductID
	}
}
func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.dataPath), 0755); err != nil {
		return err
	}
	temporary := s.dataPath + ".tmp"
	file, err := os.Create(temporary)
	if err != nil {
		return err
	}
	saved := data{s.products, s.orders, s.nextProductID}
	if err = json.NewEncoder(file).Encode(saved); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, s.dataPath)
}
func (s *Store) Products() []models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Product(nil), s.products...)
}
func (s *Store) Product(id int) (models.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.products {
		if p.ID == id {
			return p, true
		}
	}
	return models.Product{}, false
}
func (s *Store) CreateProduct(product models.Product) (models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	product.ID = s.nextProductID
	s.nextProductID++
	s.products = append(s.products, product)
	return product, s.save()
}
func (s *Store) UpdateProduct(id int, product models.Product) (models.Product, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.products {
		if s.products[i].ID == id {
			product.ID = id
			s.products[i] = product
			return product, true, s.save()
		}
	}
	return models.Product{}, false, nil
}
func (s *Store) DeleteProduct(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.products {
		if p.ID == id {
			s.products = append(s.products[:i], s.products[i+1:]...)
			return true, s.save()
		}
	}
	return false, nil
}
func (s *Store) Orders() []models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Order(nil), s.orders...)
}
func (s *Store) CreateOrder(order models.Order) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order.ID = len(s.orders) + 1
	s.orders = append(s.orders, order)
	return order, s.save()
}
func (s *Store) UpdateOrderStatus(id int, status string) (models.Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.orders {
		if s.orders[i].ID == id {
			s.orders[i].Status = status
			return s.orders[i], true, s.save()
		}
	}
	return models.Order{}, false, nil
}
