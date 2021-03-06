package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User represents the object of individual and member of organization.
type User struct {
	gorm.Model
	Email     string `gorm:"type:varchar(100);unique_index"`
	Name      string
	Password  string
	Lastlogin *time.Time
}

// CreateUser creates record of a new user.
func CreateUser(name, password, email string) (*User, error) {

	pdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	row := &User{
		Name:     name,
		Email:    email,
		Password: string(pdHash),
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		if !tx.Where("email = ?", email).First(new(User)).RecordNotFound() {
			return ErrUserAlreadyExist{0, email}
		}

		return tx.Create(row).Error
	})

	return row, err
}

// GetUser get User Info
func GetUser(email, password string) (*User, error) {
	if len(email) == 0 || len(password) == 0 {
		return nil, ErrUserNotExist{0, email}
	}

	user := new(User)
	if db.Where("email = ?", email).First(user).RecordNotFound() {
		return nil, ErrUserNotExist{0, email}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
