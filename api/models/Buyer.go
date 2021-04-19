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

type Buyer struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *Buyer) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Buyer) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *Buyer) Validate(action string) error {
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

func (u *Buyer) SaveBuyer(db *gorm.DB) (*Buyer, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Buyer{}, err
	}
	return u, nil
}

func (u *Buyer) FindAllBuyers(db *gorm.DB) (*[]Buyer, error) {
	var err error
	buyers := []Buyer{}
	err = db.Debug().Model(&Buyer{}).Limit(100).Find(&buyers).Error
	if err != nil {
		return &[]Buyer{}, err
	}
	return &buyers, err
}

func (u *Buyer) FindBuyerByID(db *gorm.DB, uid uint32) (*Buyer, error) {
	var err error
	err = db.Debug().Model(Buyer{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Buyer{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Buyer{}, errors.New("Buyer Not Found")
	}
	return u, err
}

func (u *Buyer) UpdateABuyer(db *gorm.DB, uid uint32) (*Buyer, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&Buyer{}).Where("id = ?", uid).Take(&Buyer{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"name":      u.Name,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Buyer{}, db.Error
	}
	// This is the display the updated Buyer
	err = db.Debug().Model(&Buyer{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Buyer{}, err
	}
	return u, nil
}

func (u *Buyer) DeleteABuyer(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Buyer{}).Where("id = ?", uid).Take(&Buyer{}).Delete(&Buyer{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
