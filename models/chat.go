package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Sender    string `gorm:"size:255;not null;" json:"sender"`
	Recipient string `gorm:"size:255;not null;" json:"recipient"`
	Content   string `gorm:"size:255;not null;" json:"content"`
}

func (u *Message) SaveMessage() (*Message, error) {

	err := DB.Create(&u).Error
	if err != nil {
		return &Message{}, err
	}
	return u, nil
}

func GetMessage() ([]Message, error) {

	var messages []Message

	if err := DB.Find(&messages).Error; err != nil {
		return messages, errors.New("messages not found")
	}

	return messages, nil
}
