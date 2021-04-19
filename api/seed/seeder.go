package seed

import (
	"log"

	"github.com/funthere/komodo/api/models"
	"github.com/jinzhu/gorm"
)

var buyers = []models.Buyer{
	models.Buyer{
		Name:     "Buyer 1",
		Email:    "buyer1@gmail.com",
		Password: "password",
		Address:  "alamat buyer 1",
	},
	models.Buyer{
		Name:     "Buyer 2",
		Email:    "buyer2@gmail.com",
		Password: "password",
		Address:  "alamat buyer 2",
	},
}
var sellers = []models.Seller{
	models.Seller{
		Name:     "Seller 1",
		Email:    "seller1@gmail.com",
		Password: "password",
		Address:  "alamat seller 1",
	},
	models.Seller{
		Name:     "Seller 2",
		Email:    "seller2@gmail.com",
		Password: "password",
		Address:  "alamat seller 2",
	},
}

var products = []models.Product{
	models.Product{
		ProductName: "Product Name 1",
		Description: "Hello world 1",
		Price:       100000,
		SellerID:    1,
	},
	models.Product{
		ProductName: "Produc 2",
		Description: "Hello world 2",
		Price:       150000,
		SellerID:    1,
	},
	models.Product{
		ProductName: "Product Name Seller 2",
		Description: "Desc prod ",
		Price:       50000,
		SellerID:    2,
	},
	models.Product{
		ProductName: "Produc 2 Seller 2",
		Description: "Desc prod 2",
		Price:       100000,
		SellerID:    2,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Buyer{}, &models.Product{}, &models.Seller{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Buyer{}, &models.Seller{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Product{}).AddForeignKey("seller_id", "sellers(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range buyers {
		err = db.Debug().Model(&models.Buyer{}).Create(&buyers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed buyers table: %v", err)
		}
		err = db.Debug().Model(&models.Seller{}).Create(&sellers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed sellers table: %v", err)
		}
	}

	for i, _ := range products {
		err = db.Debug().Model(&models.Product{}).Create(&products[i]).Error
		if err != nil {
			log.Fatalf("cannot seed products table: %v", err)
		}

	}
}
