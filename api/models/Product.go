package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	ProductName string    `gorm:"size:255;not null;unique" json:"product_name"`
	Description string    `gorm:"size:255;not null;" json:"description"`
	Price       uint64    `gorm:"not null" json:"price"`
	Seller      Seller    `json:"seller"`
	SellerID    uint32    `gorm:"not null" json:"seller_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.ID = 0
	p.ProductName = html.EscapeString(strings.TrimSpace(p.ProductName))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Seller = Seller{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {

	if p.ProductName == "" {
		return errors.New("Required ProductName")
	}
	if p.Description == "" {
		return errors.New("Required Description")
	}
	if p.SellerID < 1 {
		return errors.New("Required Seller")
	}
	return nil
}

func (p *Product) Save(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Seller{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAll(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i := range products {
			err := db.Debug().Model(&Seller{}).Where("id = ?", products[i].SellerID).Take(&products[i].Seller).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}

func (p *Product) FindByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Limit(100).Find(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Seller{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindBySellerID(db *gorm.DB, sid uint32) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Where("seller_id = ?", sid).Limit(100).Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i := range products {
			err := db.Debug().Model(&Seller{}).Where("id = ?", products[i].SellerID).Take(&products[i].Seller).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}

func (p *Product) Update(db *gorm.DB) (*Product, error) {

	var err error

	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{ProductName: p.ProductName, Description: p.Description, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Seller{}).Where("id = ?", p.SellerID).Take(&p.Seller).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) Delete(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and seller_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
