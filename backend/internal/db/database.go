package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"ecommerce-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductInOrders = errors.New("product is referenced by existing orders")

// Database provides persistence for products and orders.
type Database struct {
	pool *pgxpool.Pool
}

// ProductFilter narrows product queries.
type ProductFilter struct {
	Sport      string
	Gender     string
	Type       string
	ActiveOnly bool
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
	if err := database.Seed(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return database, nil
}

func (d *Database) Close() { d.pool.Close() }

// Migrate creates and updates the application schema.
func (d *Database) Migrate(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id INT8 PRIMARY KEY DEFAULT unique_rowid(),
			name STRING NOT NULL,
			description STRING NOT NULL DEFAULT '',
			sport STRING NOT NULL DEFAULT '',
			gender STRING NOT NULL DEFAULT '',
			product_type STRING NOT NULL DEFAULT '',
			price DECIMAL(12, 2) NOT NULL,
			original_price DECIMAL(12, 2),
			stock INT NOT NULL,
			image STRING NOT NULL DEFAULT '',
			image2 STRING NOT NULL DEFAULT '',
			active BOOL NOT NULL DEFAULT true
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INT8 PRIMARY KEY DEFAULT unique_rowid(),
			customer STRING NOT NULL,
			email STRING NOT NULL DEFAULT '',
			total DECIMAL(12, 2) NOT NULL,
			status STRING NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			order_id INT8 NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id INT8 NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
			quantity INT NOT NULL CHECK (quantity > 0),
			PRIMARY KEY (order_id, product_id)
		)`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS sport STRING NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS gender STRING NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS product_type STRING NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS original_price DECIMAL(12, 2)`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS image2 STRING NOT NULL DEFAULT ''`,
		`ALTER TABLE products ADD COLUMN IF NOT EXISTS category STRING NOT NULL DEFAULT ''`,
		`UPDATE products SET sport = category WHERE sport = '' AND category != ''`,
	}
	for _, stmt := range stmts {
		if _, err := d.pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("migrate failed: %w", err)
		}
	}
	return nil
}

const productColumns = `id, name, description, sport, gender, product_type, price, original_price, stock, image, image2, active`

func scanProduct(row pgx.Row) (models.Product, error) {
	var product models.Product
	err := row.Scan(
		&product.ID, &product.Name, &product.Description,
		&product.Sport, &product.Gender, &product.Type,
		&product.Price, &product.OriginalPrice, &product.Stock,
		&product.Image, &product.Image2, &product.Active,
	)
	return product, err
}

func (d *Database) Products(ctx context.Context, filter ProductFilter) ([]models.Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE 1=1`
	args := []any{}
	index := 1

	if filter.ActiveOnly {
		query += ` AND active = true`
	}
	if filter.Sport != "" {
		query += fmt.Sprintf(` AND sport = $%d`, index)
		args = append(args, filter.Sport)
		index++
	}
	if filter.Gender != "" {
		query += fmt.Sprintf(` AND gender = $%d`, index)
		args = append(args, filter.Gender)
		index++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(` AND product_type = $%d`, index)
		args = append(args, filter.Type)
	}

	query += ` ORDER BY id`

	rows, err := d.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (d *Database) Product(ctx context.Context, id int) (models.Product, bool, error) {
	row := d.pool.QueryRow(ctx, `SELECT `+productColumns+` FROM products WHERE id = $1`, id)
	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, false, nil
	}
	return product, err == nil, err
}

func (d *Database) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	err := d.pool.QueryRow(ctx,
		`INSERT INTO products (name, description, sport, gender, product_type, price, original_price, stock, image, image2, active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		product.Name, product.Description, product.Sport, product.Gender, product.Type,
		product.Price, product.OriginalPrice, product.Stock, product.Image, product.Image2, product.Active,
	).Scan(&product.ID)
	return product, err
}

func (d *Database) UpdateProduct(ctx context.Context, id int, product models.Product) (models.Product, bool, error) {
	command, err := d.pool.Exec(ctx,
		`UPDATE products SET name = $2, description = $3, sport = $4, gender = $5, product_type = $6,
		 price = $7, original_price = $8, stock = $9, image = $10, image2 = $11, active = $12
		 WHERE id = $1`,
		id, product.Name, product.Description, product.Sport, product.Gender, product.Type,
		product.Price, product.OriginalPrice, product.Stock, product.Image, product.Image2, product.Active,
	)
	if err != nil {
		return models.Product{}, false, err
	}
	if command.RowsAffected() == 0 {
		return models.Product{}, false, nil
	}
	product.ID = id
	return product, true, nil
}

func (d *Database) DeleteProduct(ctx context.Context, id int) (bool, error) {
	referenced, err := d.productInOrders(ctx, id)
	if err != nil {
		return false, err
	}
	if referenced {
		return false, ErrProductInOrders
	}
	command, err := d.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	return command.RowsAffected() > 0, err
}

func (d *Database) productInOrders(ctx context.Context, productID int) (bool, error) {
	var exists bool
	err := d.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM order_items WHERE product_id = $1)`, productID,
	).Scan(&exists)
	return exists, err
}

func (d *Database) Orders(ctx context.Context) ([]models.Order, error) {
	rows, err := d.pool.Query(ctx, `SELECT id, customer, email, total, status, created_at FROM orders ORDER BY created_at DESC`)
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
		items, err := d.orderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (d *Database) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return order, err
	}
	defer tx.Rollback(ctx)

	for _, item := range order.Items {
		var stock int
		var active bool
		err := tx.QueryRow(ctx,
			`SELECT stock, active FROM products WHERE id = $1 FOR UPDATE`, item.ProductID,
		).Scan(&stock, &active)
		if errors.Is(err, pgx.ErrNoRows) {
			return order, fmt.Errorf("product %d not found", item.ProductID)
		}
		if err != nil {
			return order, err
		}
		if !active {
			return order, fmt.Errorf("product %d is not available", item.ProductID)
		}
		if stock < item.Quantity {
			return order, fmt.Errorf("not enough stock for product %d", item.ProductID)
		}
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO orders (customer, email, total, status, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		order.Customer, order.Email, order.Total, order.Status, order.CreatedAt,
	).Scan(&order.ID)
	if err != nil {
		return order, err
	}

	for _, item := range order.Items {
		if _, err := tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			order.ID, item.ProductID, item.Quantity,
		); err != nil {
			return order, err
		}
		if _, err := tx.Exec(ctx,
			`UPDATE products SET stock = stock - $2 WHERE id = $1`, item.ProductID, item.Quantity,
		); err != nil {
			return order, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return order, err
	}
	return order, nil
}

func (d *Database) UpdateOrderStatus(ctx context.Context, id int, status string) (models.Order, bool, error) {
	var order models.Order
	err := d.pool.QueryRow(ctx,
		`UPDATE orders SET status = $2 WHERE id = $1 RETURNING id, customer, email, total, status, created_at`,
		id, status,
	).Scan(&order.ID, &order.Customer, &order.Email, &order.Total, &order.Status, &order.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, false, nil
	}
	if err != nil {
		return models.Order{}, false, err
	}
	order.Items, err = d.orderItems(ctx, order.ID)
	return order, err == nil, err
}

func (d *Database) orderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	rows, err := d.pool.Query(ctx,
		`SELECT product_id, quantity FROM order_items WHERE order_id = $1 ORDER BY product_id`, orderID,
	)
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

// NormalizeFilter converts query params into a ProductFilter.
func NormalizeFilter(sport, gender, productType string) ProductFilter {
	return ProductFilter{
		Sport:  strings.TrimSpace(sport),
		Gender: strings.TrimSpace(gender),
		Type:   strings.TrimSpace(productType),
	}
}

// WithTimeout wraps a background context with a default timeout.
func WithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
