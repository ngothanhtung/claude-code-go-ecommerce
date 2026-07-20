package catalog

import (
	"time"

	"github.com/lib/pq"
)

type Category struct {
	ID        string    `gorm:"primaryKey;size:50" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	IconName  string    `gorm:"size:100" json:"icon_name"`
	ColorHex  string    `gorm:"size:9" json:"color_hex"`
	CreatedAt time.Time `json:"created_at"`
}

func (Category) TableName() string { return "categories" }

type Product struct {
	ID           string         `gorm:"primaryKey;size:100" json:"id"`
	Name         string         `gorm:"size:255" json:"name"`
	Price        float64        `gorm:"type:decimal(10,2)" json:"price"`
	Rating       float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`
	ReviewsCount int            `gorm:"column:reviews_count;default:0" json:"reviews_count"`
	IconName     string         `gorm:"size:100" json:"icon_name"`
	ColorHex     string         `gorm:"size:9" json:"color_hex"`
	CategoryID   string         `gorm:"size:50" json:"category_id"`
	Description  string         `gorm:"type:text" json:"description"`
	ImageURLs    pq.StringArray `gorm:"type:text[]" json:"image_urls"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index" json:"-"`
}

func (Product) TableName() string { return "products" }

type ProductPublic struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Rating       float64   `json:"rating"`
	ReviewsCount int       `json:"reviews_count"`
	IconName     string    `json:"icon_name"`
	ColorHex     string    `json:"color_hex"`
	CategoryID   string    `json:"category_id"`
	Description  string    `json:"description"`
	ImageURLs    []string  `json:"image_urls"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p Product) ToPublic() ProductPublic {
	urls := []string(p.ImageURLs)
	if urls == nil {
		urls = []string{}
	}
	return ProductPublic{
		ID: p.ID, Name: p.Name, Price: p.Price, Rating: p.Rating,
		ReviewsCount: p.ReviewsCount, IconName: p.IconName, ColorHex: p.ColorHex,
		CategoryID: p.CategoryID, Description: p.Description,
		ImageURLs: urls, CreatedAt: p.CreatedAt,
	}
}

type Promo struct {
	ID        string    `gorm:"primaryKey;size:100" json:"id"`
	Badge     string    `gorm:"size:100" json:"badge"`
	Title     string    `gorm:"size:255" json:"title"`
	Subtitle  string    `gorm:"size:255" json:"subtitle"`
	IconName  string    `gorm:"size:100" json:"icon_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (Promo) TableName() string { return "promos" }
