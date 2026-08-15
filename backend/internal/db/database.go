package db

import (
	"context"
	"errors"
	"os"
	"time"

	"ecommerce-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database provides persistence for products and orders.
type Database struct {
	pool *pgxpool.Pool
}

// Connect opens a CockroachDB/PostgreSQL connection using DATABASE_URL.
func Connect(ctx context.Context) (*Database, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	database := &Database{pool: pool}
	if err := database.Migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return database, nil
}

func (d *Database) Close() { d.pool.Close() }

// Migrate creates the application's initial schema if it does not already exist.
func (d *Database) Migrate(ctx context.Context) error {
	_, err := d.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id INT8 PRIMARY KEY DEFAULT unique_rowid(),
			name STRING NOT NULL,
			description STRING NOT NULL DEFAULT '',
			category STRING NOT NULL DEFAULT '',
			price DECIMAL(12, 2) NOT NULL,
			stock INT NOT NULL,
			image STRING NOT NULL DEFAULT '',
			active BOOL NOT NULL DEFAULT true
		);
		CREATE TABLE IF NOT EXISTS orders (
			id INT8 PRIMARY KEY DEFAULT unique_rowid(),
			customer STRING NOT NULL,
			email STRING NOT NULL DEFAULT '',
			total DECIMAL(12, 2) NOT NULL,
			status STRING NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
		CREATE TABLE IF NOT EXISTS order_items (
			order_id INT8 NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id INT8 NOT NULL REFERENCES products(id),
			quantity INT NOT NULL CHECK (quantity > 0),
			PRIMARY KEY (order_id, product_id)
		);`)
	return err
}

func (d *Database) Products() ([]models.Product, error) {
	rows, err := d.pool.Query(context.Background(), `SELECT id, name, description, category, price, stock, image, active FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Stock, &product.Image, &product.Active); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (d *Database) Product(id int) (models.Product, bool, error) {
	var product models.Product
	err := d.pool.QueryRow(context.Background(), `SELECT id, name, description, category, price, stock, image, active FROM products WHERE id = $1`, id).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.Stock, &product.Image, &product.Active)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, false, nil
	}
	return product, err == nil, err
}

func (d *Database) CreateProduct(product models.Product) (models.Product, error) {
	err := d.pool.QueryRow(context.Background(), `INSERT INTO products (name, description, category, price, stock, image, active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, product.Name, product.Description, product.Category, product.Price, product.Stock, product.Image, product.Active).Scan(&product.ID)
	return product, err
}

func (d *Database) UpdateProduct(id int, product models.Product) (models.Product, bool, error) {
	command, err := d.pool.Exec(context.Background(), `UPDATE products SET name = $2, description = $3, category = $4, price = $5, stock = $6, image = $7, active = $8 WHERE id = $1`, id, product.Name, product.Description, product.Category, product.Price, product.Stock, product.Image, product.Active)
	if err != nil {
		return models.Product{}, false, err
	}
	if command.RowsAffected() == 0 {
		return models.Product{}, false, nil
	}
	product.ID = id
	return product, true, nil
}

func (d *Database) DeleteProduct(id int) (bool, error) {
	command, err := d.pool.Exec(context.Background(), `DELETE FROM products WHERE id = $1`, id)
	return command.RowsAffected() > 0, err
}

func (d *Database) Orders() ([]models.Order, error) {
	rows, err := d.pool.Query(context.Background(), `SELECT id, customer, email, total, status, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.Customer, &order.Email, &order.Total, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		items, err := d.orderItems(context.Background(), order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (d *Database) CreateOrder(order models.Order) (models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return order, err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `INSERT INTO orders (customer, email, total, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, order.Customer, order.Email, order.Total, order.Status, order.CreatedAt).Scan(&order.ID)
	if err != nil {
		return order, err
	}
	for _, item := range order.Items {
		if _, err := tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`, order.ID, item.ProductID, item.Quantity); err != nil {
			return order, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return order, err
	}
	return order, nil
}

func (d *Database) UpdateOrderStatus(id int, status string) (models.Order, bool, error) {
	var order models.Order
	err := d.pool.QueryRow(context.Background(), `UPDATE orders SET status = $2 WHERE id = $1 RETURNING id, customer, email, total, status, created_at`, id, status).Scan(&order.ID, &order.Customer, &order.Email, &order.Total, &order.Status, &order.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, false, nil
	}
	if err != nil {
		return models.Order{}, false, err
	}
	order.Items, err = d.orderItems(context.Background(), order.ID)
	return order, err == nil, err
}

func (d *Database) orderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := d.pool.Query(ctx, `SELECT product_id, quantity FROM order_items WHERE order_id = $1 ORDER BY product_id`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
