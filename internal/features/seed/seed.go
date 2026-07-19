package seed

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ngothanhtung/go-tutorials/internal/features/rbac"
	usermodel "github.com/ngothanhtung/go-tutorials/internal/features/user"
)

// Run creates default roles and an admin user if they do not exist.
func Run(ctx context.Context, db *gorm.DB, adminEmail, adminPassword string) error {
	adminRole := rbac.Role{Name: "admin"}
	if err := db.WithContext(ctx).Where("name = ?", "admin").FirstOrCreate(&adminRole).Error; err != nil {
		return err
	}
	userRole := rbac.Role{Name: "user"}
	if err := db.WithContext(ctx).Where("name = ?", "user").FirstOrCreate(&userRole).Error; err != nil {
		return err
	}

	var count int64
	if err := db.WithContext(ctx).Model(&usermodel.User{}).Where("email = ?", adminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("seed: admin already exists, skipping")
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	rid := adminRole.ID
	admin := usermodel.User{
		ID:           uuid.New(),
		Email:        adminEmail,
		Name:         "Administrator",
		PasswordHash: string(hash),
		RoleID:       &rid,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := db.WithContext(ctx).Create(&admin).Error; err != nil {
		return err
	}
	log.Println("seed: admin user created:", adminEmail)
	return nil
}
