package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Merchant struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	User         User      `json:"users"`
	UserID       uint32    `gorm:"not null" json:"user_id"`
	MerchantName string    `gorm:"size:40;not null;unique" json:"merchant_name"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// func (m *Merchant) FindAllMerchant(db *gorm.DB) (*[]Merchant, error) {
// 	var err error
// 	merchants := []Merchant{}
// 	err = db.Debug().Model(&Merchant{}).Limit(100).Find(&merchants).Error
// 	if err != nil {
// 		return &[]Merchant{}, err
// 	}
// 	if len(merchants) > 0 {
// 		for i := range merchants {
// 			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
// 			if err != nil {
// 				return &[]Post{}, err
// 			}
// 		}
// 	}
// 	return &merchants, nil
// }

func (u *Merchant) GetAllMerchant(db *gorm.DB) (*[]Merchant, error) {
	var err error
	merchants := []Merchant{}
	err = db.Debug().Model(&Merchant{}).Limit(100).Find(&merchants).Error
	if err != nil {
		return &[]Merchant{}, err
	}

	return &merchants, err
}
