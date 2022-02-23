package seed

import (
	"log"
	"math/rand"

	"github.com/abdulhamidnugroho/go-full/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Name:     "Abdul",
		Username: "abdul",
		Password: "password",
	},
	{
		Name:     "Hamid",
		Username: "hamid",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello #1",
	},
	{
		Title:   "Title 2",
		Content: "Hello #2",
	},
}

var merchants = []models.Merchant{
	{
		MerchantName: "Merchant X",
	},
	{
		MerchantName: "Merchant Y",
	},
}

var outlets = []models.Outlet{
	{
		OutletName: "Outlet P",
	},
	{
		OutletName: "Outlet Q",
	},
}

var transactions = []models.Transaction{
	{
		BillTotal: 205800,
	},
	{
		BillTotal: 165500,
	},
	{
		BillTotal: 302800,
	},
	{
		BillTotal: 135500,
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.Transaction{}, &models.Outlet{}, &models.Merchant{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Merchant{}, &models.Outlet{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Merchant{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Outlet{}).AddForeignKey("merchant_id", "merchants(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("merchant_id", "merchants(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Transaction{}).AddForeignKey("outlet_id", "outlets(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID
		merchants[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		err = db.Debug().Model(&models.Merchant{}).Create(&merchants[i]).Error
		if err != nil {
			log.Fatalf("cannot seed merchants table: %v", err)
		}
	}

	for i := range merchants {
		outlets[i].MerchantID = merchants[i].ID

		err = db.Debug().Model(&models.Outlet{}).Create(&outlets[i]).Error
		if err != nil {
			log.Fatalf("cannot seed outlets table: %v", err)
		}
	}

	for i := range transactions {
		outlet := outlets[rand.Intn(len(outlets))]

		transactions[i].MerchantID = outlet.MerchantID
		transactions[i].OutletID = outlet.ID

		err = db.Debug().Model(&models.Transaction{}).Create(&transactions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed outlets table: %v", err)
		}
	}
}
