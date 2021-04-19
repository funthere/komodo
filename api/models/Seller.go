package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Seller struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func HashSeller(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPasswordSeller(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *Seller) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Seller) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *Seller) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *Seller) Save(db *gorm.DB) (*Seller, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Seller{}, err
	}
	return u, nil
}

func (u *Seller) FindAll(db *gorm.DB) (*[]Seller, error) {
	var err error
	sellers := []Seller{}
	err = db.Debug().Model(&Seller{}).Limit(100).Find(&sellers).Error
	if err != nil {
		return &[]Seller{}, err
	}
	return &sellers, err
}

func (u *Seller) FindByID(db *gorm.DB, uid uint32) (*Seller, error) {
	var err error
	err = db.Debug().Model(Seller{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Seller{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Seller{}, errors.New("Seller Not Found")
	}
	return u, err
}

func (u *Seller) Update(db *gorm.DB, uid uint32) (*Seller, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Seller{}).Where("id = ?", uid).Take(&Seller{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"name":      u.Name,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Seller{}, db.Error
	}
	// This is the display the updated Seller
	err = db.Debug().Model(&Seller{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Seller{}, err
	}
	return u, nil
}

func (u *Seller) Delete(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Seller{}).Where("id = ?", uid).Take(&Seller{}).Delete(&Seller{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
