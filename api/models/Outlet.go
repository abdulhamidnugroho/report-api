package models

import (
	"time"
)

type Outlet struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Merchant   Merchant  `json:"merchant"`
	MerchantID uint32    `gorm:"not null" json:"merchant_id"`
	OutletName string    `gorm:"size:40;not null;" json:"outlet_name"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
