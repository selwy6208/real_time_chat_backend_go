package models

import (
	"github.com/jinzhu/gorm"
)

type SocketMessage struct {
	gorm.Model
	MessageType string `gorm:"not null;" json:"message_type"`
	MessageData string `gorm:"not null;" json:"message_data"`
}
type IncomingNewMessage struct {
	gorm.Model
	MessageType string `gorm:"not null;" json:"message_type"`
	MessageData string `gorm:"not null;" json:"message_data"`
}

type Message struct {
	gorm.Model
	Content   string `gorm:"not null;" json:"content"`
	Sender    string `gorm:"not null;" json:"sender"`
	Recipient string `gorm:"not null;" json:"recipient"`
}

type InputMessage struct {
	gorm.Model
	Sender    string `gorm:"size:255;not null;" json:"sender"`
	Recipient string `gorm:"size:255;not null;" json:"recipient"`
}

func (u *Message) SaveMessage() (*Message, error) {
	err := DB.Create(&u).Error
	if err != nil {
		return &Message{}, err
	}
	return u, nil
}

func GetMessage(u *InputMessage) (*Message, error) {

	message := &Message{
		Model: gorm.Model{},
	}
	err := DB.Where("sender = ? AND recipient = ?", u.Sender, u.Recipient).Find(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}
