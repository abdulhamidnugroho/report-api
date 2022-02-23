package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Merchant   Merchant  `json:"merchant"`
	MerchantID uint32    `gorm:"not null;" json:"merchant_id"`
	Outlet     Outlet    `json:"outlet"`
	OutletID   uint32    `gorm:"not null;" json:"outlet_id"`
	BillTotal  float32   `gorm:"type:float;not null;" json:"bill_total"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *Transaction) GetAllTransaction(db *gorm.DB) (*[]Transaction, error) {
	var err error
	transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}

	return &transactions, err
}
