package models

import (
	"errors"
	"html"
	"strings"

	"real-chat-backend/utils/token"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"size:255;not null" json:"firstname"`
	LastName  string `gorm:"size:255;not null" json:"lastname"`
	Email     string `gorm:"size:255;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null;" json:"password"`
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found")
	}

	u.PrepareGive()

	return u, nil
}

func GetMessagesByUserID(myID uint, chatUserId uint) ([]Message, error) {
	var messages []Message

	if err := DB.Model(Message{}).Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)", myID, chatUserId, chatUserId, myID).Find(&messages).Error; err != nil {
		return messages, errors.New("users not found")
	}

	return messages, nil
}

func GetUsers() ([]User, error) {

	var users []User

	if err := DB.Find(&users).Error; err != nil {
		return users, errors.New("users not found")
	}

	return users, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email string, password string) (string, error) {

	var err error

	u := User{}

	err = DB.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser() (*User, error) {

	err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))

	return nil

}
