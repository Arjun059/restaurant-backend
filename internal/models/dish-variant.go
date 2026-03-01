package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type DishVariant struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DishID      uuid.UUID      `json:"dishId" gorm:"type:uuid;index"`
	Dish        Dish           `json:"dish" gorm:"foreignKey:DishID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	
	Name        string         `json:"name" validate:"required,min=2,max=50"`
	Price       float64        `json:"price" validate:"required,gt=0"`
	DisplayPrice float64       `json:"displayPrice"`
	
	Serving     string         `json:"serving"`     // e.g., "1 person", "2 people"
	PreparationTime string     `json:"preparationTime"`
	
	IsAvailable bool           `json:"isAvailable" gorm:"default:true"`
	
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
