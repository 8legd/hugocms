package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name  string
	Image UserImage
}

func (user *User) DisplayName() string {
	return user.Name
}

func (u *User) AfterCreate() (err error) {
	if false {
		err = errors.New("User Create Error")
	} else {
		fmt.Println("Created User", u)
	}
	return
}

func (u *User) AfterUpdate() (err error) {
	if false {
		err = errors.New("User Update Error")
	} else {
		fmt.Println("Updated User", u)
	}
	return
}

func (u *User) AfterSave() (err error) {
	if false {
		err = errors.New("User Save Error")
	} else {
		fmt.Println("Saved User", u)
	}
	return
}

func (u *User) AfterDelete() (err error) {
	if false {
		err = errors.New("User Delete Error")
	} else {
		fmt.Println("Deleted User", u)
	}
	return
}

type UserImage struct {
	gorm.Model
	UserID uint
	Image  SimpleImageStorage `sql:"type:varchar(4096)"`
}
