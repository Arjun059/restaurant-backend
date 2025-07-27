package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
	RestaurantID  uint   `json:"restaurantId"`
	Restaurant    *Restaurant  `gorm:"foreignKey:RestaurantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"` // Automatically managed by GORM
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"` // Automatically managed by GORM

	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // Enables soft delete
}
