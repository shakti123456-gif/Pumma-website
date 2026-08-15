package models

type Product struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Sport         string   `json:"sport"`
	Gender        string   `json:"gender"`
	Type          string   `json:"type"`
	Price         float64  `json:"price"`
	OriginalPrice *float64 `json:"originalPrice,omitempty"`
	Stock         int      `json:"stock"`
	Image         string   `json:"image"`
	Image2        string   `json:"image2"`
	Active        bool     `json:"active"`
}
