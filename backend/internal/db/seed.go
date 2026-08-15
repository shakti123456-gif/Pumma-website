package db

import (
	"context"

	"ecommerce-backend/internal/models"
)

type seedProduct struct {
	Name          string
	Description   string
	Sport         string
	Gender        string
	Type          string
	Price         float64
	OriginalPrice *float64
	Stock         int
	Image         string
	Image2        string
}

func (d *Database) Seed(ctx context.Context) error {
	var count int
	if err := d.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	pics := [][2]string{
		{"photo-1552346154-21d32810aba3", "photo-1542291026-7eec264c27ff"},
		{"photo-1608231387042-66d1773070a5", "photo-1600185365483-26d7a4cc7519"},
		{"photo-1575537302964-96cd47c06b1b", "photo-1600185365926-3a2ce3cdb9eb"},
		{"photo-1606107557195-0e29a4b5b4aa", "photo-1603808033192-082d6919d3e1"},
		{"photo-1605348532760-6753d2c43329", "photo-1460353581641-37baddab0fa2"},
		{"photo-1517836357463-d25dfeac3438", "photo-1518611012118-696072aa579a"},
		{"photo-1538805060514-97d9cc17730c", "photo-1523398002811-999ca8dec234"},
		{"photo-1515886657613-9f3515b0c78f", "photo-1558618666-fcd25c85cd64"},
	}

	catalog := []seedProduct{
		{"Velocity NXT", "Lightweight running shoe for daily miles.", "Running", "Men", "Shoes", 125, nil, 24, "", ""},
		{"RS-01 Runner", "Responsive cushioning for long runs.", "Running", "Women", "Shoes", 110, floatPtr(145), 18, "", ""},
		{"Sprint Elite", "Breathable running top.", "Running", "Men", "Clothing", 58, nil, 30, "", ""},
		{"Cloud Knit Tee", "Soft knit tee for easy pace days.", "Running", "Women", "Clothing", 45, nil, 28, "", ""},
		{"Dash Runner Short", "Lightweight shorts with inner brief.", "Running", "Men", "Clothing", 42, nil, 35, "", ""},
		{"Arc Runner Legging", "High-rise legging with phone pocket.", "Running", "Women", "Clothing", 68, nil, 22, "", ""},
		{"Lift Training Tee", "Moisture-wicking gym tee.", "Training", "Men", "Clothing", 38, nil, 40, "", ""},
		{"Reform Leggings", "Compression leggings for studio sessions.", "Training", "Women", "Clothing", 72, nil, 26, "", ""},
		{"Core Half-Zip", "Layer-ready half zip for warm-ups.", "Training", "Men", "Clothing", 65, nil, 20, "", ""},
		{"Zone Sports Bra", "Medium support sports bra.", "Training", "Women", "Clothing", 48, nil, 32, "", ""},
		{"Pro Training Bag", "Durable duffel for gym and travel.", "Training", "Unisex", "Clothing", 85, nil, 15, "", ""},
		{"Studio Seamless Top", "Seamless crop for training.", "Training", "Women", "Clothing", 55, floatPtr(90), 19, "", ""},
		{"Striker Pro Boot", "Firm-ground football boot.", "Football", "Men", "Shoes", 140, nil, 16, "", ""},
		{"Pitch Control Boot", "Agile touch boot for quick cuts.", "Football", "Women", "Shoes", 130, nil, 14, "", ""},
		{"Matchday Jersey", "Lightweight match jersey.", "Football", "Men", "Clothing", 75, nil, 25, "", ""},
		{"Elite Football Short", "Ventilated football short.", "Football", "Women", "Clothing", 52, nil, 21, "", ""},
		{"Rally Court", "Court shoe with stable base.", "Basketball", "Men", "Shoes", 115, nil, 17, "", ""},
		{"Lumen Court", "Low-profile women's court shoe.", "Basketball", "Women", "Shoes", 105, nil, 13, "", ""},
		{"Hoops Mesh Short", "Classic mesh basketball short.", "Basketball", "Men", "Clothing", 48, nil, 29, "", ""},
		{"Court Flex Tank", "Flexible tank for court movement.", "Basketball", "Women", "Clothing", 44, nil, 27, "", ""},
		{"Form Classic", "Everyday lifestyle tee.", "Lifestyle", "Men", "Clothing", 78, nil, 33, "", ""},
		{"Nova Luxe Hoodie", "Premium fleece hoodie.", "Lifestyle", "Women", "Clothing", 95, nil, 24, "", ""},
		{"Everyday Cargo", "Relaxed cargo pant.", "Lifestyle", "Men", "Clothing", 88, nil, 20, "", ""},
		{"Futuro Hoodie", "Oversized lifestyle hoodie.", "Lifestyle", "Women", "Clothing", 82, floatPtr(117), 18, "", ""},
		{"Drift Suede", "Minimal lifestyle sneaker.", "Lifestyle", "Men", "Shoes", 120, nil, 12, "", ""},
		{"Vibe Slipstream", "Clean everyday sneaker.", "Lifestyle", "Women", "Shoes", 98, nil, 16, "", ""},
		{"Orbit Windbreaker", "Packable windbreaker shell.", "Lifestyle", "Men", "Clothing", 140, nil, 11, "", ""},
		{"Pulse Rib Top", "Ribbed long-sleeve top.", "Lifestyle", "Women", "Clothing", 62, nil, 23, "", ""},
	}

	for i, item := range catalog {
		pair := pics[i%len(pics)]
		item.Image = unsplash(pair[0])
		item.Image2 = unsplash(pair[1])
		if _, err := d.CreateProduct(ctx, models.Product{
			Name:          item.Name,
			Description:   item.Description,
			Sport:         item.Sport,
			Gender:        item.Gender,
			Type:          item.Type,
			Price:         item.Price,
			OriginalPrice: item.OriginalPrice,
			Stock:         item.Stock,
			Image:         item.Image,
			Image2:        item.Image2,
			Active:        true,
		}); err != nil {
			return err
		}
	}
	return nil
}

func unsplash(id string) string {
	return "https://images.unsplash.com/" + id + "?auto=format&fit=crop&w=900&q=85"
}

func floatPtr(value float64) *float64 {
	return &value
}
