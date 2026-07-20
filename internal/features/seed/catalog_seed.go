package seed

import (
	"context"

	"gorm.io/gorm"

	"github.com/ngothanhtung/go-tutorials/internal/features/catalog"
)

func SeedCatalog(ctx context.Context, db *gorm.DB) error {
	// Categories
	categories := []catalog.Category{
		{ID: "fashion", Name: "Fashion", IconName: "checkroom", ColorHex: "FFE91E63"},
		{ID: "electronics", Name: "Electronics", IconName: "devices", ColorHex: "FF3F51B5"},
		{ID: "home", Name: "Home", IconName: "home", ColorHex: "FFFF5722"},
		{ID: "beauty", Name: "Beauty", IconName: "spa", ColorHex: "FFE91E63"},
		{ID: "sports", Name: "Sports", IconName: "sports_soccer", ColorHex: "FF4CAF50"},
		{ID: "travel", Name: "Travel", IconName: "flight", ColorHex: "FF2196F3"},
	}
	for _, c := range categories {
		if err := db.WithContext(ctx).
			Where("id = ?", c.ID).
			FirstOrCreate(&c).Error; err != nil {
			return err
		}
	}

	// Products — use the exact 24 products from Flutter seed_data.dart.
	// IMPORTANT: Use EXACT product IDs p1-p24 and EXACT category IDs.
	// image_urls from Flutter product.imageUrls field (or generated picsum if empty).
	// Copy the full 24 product entries exactly. Do NOT truncate.
	products := []catalog.Product{
		{ID: "p1", Name: "Linen Oversized Shirt", Price: 59, Rating: 4.2, ReviewsCount: 23, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "A breezy linen shirt for warm days.", ImageURLs: []string{}},
		{ID: "p2", Name: "Linen Wide Trousers", Price: 49, Rating: 4.1, ReviewsCount: 18, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Comfortable wide-leg trousers.", ImageURLs: []string{}},
		{ID: "p3", Name: "Linen Short-Sleeve Shirt", Price: 54, Rating: 3.8, ReviewsCount: 12, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Relaxed short-sleeve linen.", ImageURLs: []string{}},
		{ID: "p4", Name: "Unstructured Linen Jacket", Price: 89, Rating: 4.7, ReviewsCount: 45, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Lightweight unstructured jacket.", ImageURLs: []string{}},
		{ID: "p5", Name: "Lightweight Puffer Jacket", Price: 109, Rating: 4.5, ReviewsCount: 34, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Warm but lightweight puffer jacket.", ImageURLs: []string{}},
		{ID: "p6", Name: "Sleeveless Tank Top", Price: 24, Rating: 3.9, ReviewsCount: 8, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Breathable sleeveless top.", ImageURLs: []string{}},
		{ID: "p7", Name: "Classic Chino Pants", Price: 69, Rating: 4.3, ReviewsCount: 28, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Timeless chino pants.", ImageURLs: []string{}},
		{ID: "p8", Name: "Striped Rugby Shirt", Price: 64, Rating: 4.0, ReviewsCount: 16, IconName: "checkroom", ColorHex: "FFE91E63", CategoryID: "fashion", Description: "Classic striped rugby shirt.", ImageURLs: []string{}},
		{ID: "p9", Name: "Wireless Earbuds", Price: 79, Rating: 4.6, ReviewsCount: 89, IconName: "devices", ColorHex: "FF3F51B5", CategoryID: "electronics", Description: "Premium sound, no wires.", ImageURLs: []string{}},
		{ID: "p10", Name: "Staub Cast Iron Dutch Oven", Price: 219, Rating: 4.8, ReviewsCount: 67, IconName: "home", ColorHex: "FFFF5722", CategoryID: "home", Description: "Enameled cast iron for perfect braising.", ImageURLs: []string{}},
		{ID: "p11", Name: "Vitamix A3500i", Price: 549, Rating: 4.9, ReviewsCount: 43, IconName: "home", ColorHex: "FFFF5722", CategoryID: "home", Description: "Professional-grade blender.", ImageURLs: []string{}},
		{ID: "p12", Name: "Air Purifier", Price: 149, Rating: 4.2, ReviewsCount: 31, IconName: "home", ColorHex: "FFFF5722", CategoryID: "home", Description: "HEPA filtration for cleaner air.", ImageURLs: []string{}},
		{ID: "p13", Name: "Down Pillow", Price: 49, Rating: 4.4, ReviewsCount: 52, IconName: "home", ColorHex: "FFFF5722", CategoryID: "home", Description: "Cloud-soft down pillow.", ImageURLs: []string{}},
		{ID: "p14", Name: "Smart Air Monitor", Price: 39, Rating: 3.7, ReviewsCount: 19, IconName: "devices", ColorHex: "FF3F51B5", CategoryID: "electronics", Description: "Track your home air quality.", ImageURLs: []string{}},
		{ID: "p15", Name: "Ceramic Plant Pot", Price: 19, Rating: 4.1, ReviewsCount: 24, IconName: "home", ColorHex: "FFFF5722", CategoryID: "home", Description: "Minimalist ceramic planter.", ImageURLs: []string{}},
		{ID: "p16", Name: "Vitamin C Serum", Price: 38, Rating: 4.7, ReviewsCount: 93, IconName: "spa", ColorHex: "FFE91E63", CategoryID: "beauty", Description: "Brightening daily vitamin C serum.", ImageURLs: []string{}},
		{ID: "p17", Name: "Retinol Night Cream", Price: 45, Rating: 4.5, ReviewsCount: 71, IconName: "spa", ColorHex: "FFE91E63", CategoryID: "beauty", Description: "Anti-aging retinol night cream.", ImageURLs: []string{}},
		{ID: "p18", Name: "Hyaluronic Acid Moisturizer", Price: 42, Rating: 4.6, ReviewsCount: 82, IconName: "spa", ColorHex: "FFE91E63", CategoryID: "beauty", Description: "Deep hydration hyaluronic moisturizer.", ImageURLs: []string{}},
		{ID: "p19", Name: "Insulated Water Bottle", Price: 29, Rating: 4.8, ReviewsCount: 134, IconName: "sports_soccer", ColorHex: "FF4CAF50", CategoryID: "sports", Description: "Keep drinks cold for 24 hours.", ImageURLs: []string{}},
		{ID: "p20", Name: "Travel Neck Pillow", Price: 19, Rating: 4.3, ReviewsCount: 58, IconName: "flight", ColorHex: "FF2196F3", CategoryID: "travel", Description: "Memory foam travel pillow.", ImageURLs: []string{}},
		{ID: "p21", Name: "Resistance Band Set", Price: 24, Rating: 4.2, ReviewsCount: 47, IconName: "sports_soccer", ColorHex: "FF4CAF50", CategoryID: "sports", Description: "5-level resistance band set.", ImageURLs: []string{}},
		{ID: "p22", Name: "Foam Roller", Price: 34, Rating: 4.4, ReviewsCount: 39, IconName: "sports_soccer", ColorHex: "FF4CAF50", CategoryID: "sports", Description: "High-density muscle roller.", ImageURLs: []string{}},
		{ID: "p23", Name: "Yoga Mat", Price: 49, Rating: 4.9, ReviewsCount: 112, IconName: "sports_soccer", ColorHex: "FF4CAF50", CategoryID: "sports", Description: "Non-slip eco-friendly yoga mat.", ImageURLs: []string{}},
		{ID: "p24", Name: "Yoga Mat", Price: 49, Rating: 4.9, ReviewsCount: 112, IconName: "sports_soccer", ColorHex: "FF4CAF50", CategoryID: "sports", Description: "Non-slip eco-friendly yoga mat.", ImageURLs: []string{}},
	}
	for _, p := range products {
		if err := db.WithContext(ctx).
			Where("id = ?", p.ID).
			FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}

	// Promos
	promos := []catalog.Promo{
		{ID: "new", Badge: "NEW", Title: "New Arrivals", Subtitle: "Check out the latest drops", IconName: "local_offer"},
		{ID: "sale", Badge: "SALE", Title: "Summer Sale", Subtitle: "Up to 50% off selected items", IconName: "sell"},
		{ID: "trending", Badge: "TRENDING", Title: "Trending Now", Subtitle: "Most popular picks this week", IconName: "trending_up"},
	}
	for _, p := range promos {
		if err := db.WithContext(ctx).
			Where("id = ?", p.ID).
			FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	return nil
}
