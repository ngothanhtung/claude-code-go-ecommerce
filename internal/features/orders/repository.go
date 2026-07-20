package orders

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Repository interface {
	ListForUser(ctx context.Context, userID uuid.UUID) ([]Order, error)
	GetForUser(ctx context.Context, userID, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*Order, error)
	// ClearCart removes all cart_items rows for the user (called during checkout).
	ClearCart(ctx context.Context, userID uuid.UUID) error
	// productLookup returns (name, price) for a product id.
	ProductLookup(ctx context.Context, productID string) (name string, price float64, err error)
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) ListForUser(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	var rows []Order
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, apperr.NewInternal("list orders", err)
	}
	return rows, nil
}

func (r *repo) GetForUser(ctx context.Context, userID, id uuid.UUID) (*Order, error) {
	var row Order
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, id).
		Preload("Items").
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("order not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("get order", err)
	}
	return &row, nil
}

func (r *repo) ProductLookup(ctx context.Context, productID string) (string, float64, error) {
	type prodRow struct {
		Name  string
		Price float64
	}
	var p prodRow
	err := r.db.WithContext(ctx).
		Table("products").
		Select("name, price").
		Where("id = ?", productID).
		Scan(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", 0, apperr.NewNotFound("product not found")
	}
	if err != nil {
		return "", 0, apperr.NewInternal("lookup product", err)
	}
	if p.Name == "" {
		return "", 0, apperr.NewNotFound("product not found")
	}
	return p.Name, p.Price, nil
}

func (r *repo) Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*Order, error) {
	var created *Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Snapshot each line — fetch product name/price server-side and
		//    compute the total. If any product is missing, the txn rolls back.
		snapshots := make([]OrderItem, 0, len(req.Items))
		var total float64
		for _, it := range req.Items {
			type prodRow struct {
				Name  string
				Price float64
			}
			var p prodRow
			txErr := tx.WithContext(ctx).
				Table("products").
				Select("name, price").
				Where("id = ?", it.ProductID).
				Scan(&p).Error
			if txErr != nil {
				return apperr.NewInternal("lookup product", txErr)
			}
			if p.Name == "" {
				return apperr.NewNotFound("product not found")
			}
			snapshots = append(snapshots, OrderItem{
				ProductID: it.ProductID,
				Name:      p.Name,
				UnitPrice: p.Price,
				Quantity:  it.Quantity,
			})
			total += p.Price * float64(it.Quantity)
		}

		// 2) Insert header row. Status defaults to "processing".
		o := &Order{
			UserID:          userID,
			Total:           total,
			Status:          "processing",
			ShippingAddress: req.ShippingAddress,
			PaymentMethod:   req.PaymentMethod,
			Items:           snapshots,
		}
		// Items are inserted explicitly below after the generated order ID is
		// available; omit the association here to avoid inserting them twice.
		if err := tx.Omit("Items").Create(o).Error; err != nil {
			return apperr.NewInternal("create order", err)
		}

		// 3) Wire snapshots to the order id and insert them.
		for i := range snapshots {
			snapshots[i].OrderID = o.ID
		}
		if err := tx.Create(&snapshots).Error; err != nil {
			if strings.Contains(err.Error(), "violates foreign key") {
				return apperr.NewNotFound("order reference missing")
			}
			return apperr.NewInternal("insert order items", err)
		}
		o.Items = snapshots

		// 4) Clear the user's cart after a successful checkout.
		if err := tx.Exec("DELETE FROM cart_items WHERE user_id = ?", userID).Error; err != nil {
			return apperr.NewInternal("clear cart", err)
		}

		created = o
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *repo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*Order, error) {
	if err := r.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return nil, apperr.NewInternal("update order status", err)
	}
	var row Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&row, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("order not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("reload order", err)
	}
	return &row, nil
}

func (r *repo) ClearCart(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Exec("DELETE FROM cart_items WHERE user_id = ?", userID).Error; err != nil {
		return apperr.NewInternal("clear cart", err)
	}
	return nil
}
