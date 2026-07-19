package rbac

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string       `gorm:"uniqueIndex;size:50" json:"name"`
	CreatedAt   time.Time    `json:"created_at"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:100" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName ensures correct mapping for the role entity.
func (Role) TableName() string { return "roles" }

var _ = gorm.Model{}
