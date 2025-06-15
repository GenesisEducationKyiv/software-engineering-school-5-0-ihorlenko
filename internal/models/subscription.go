package models

import "time"

type Subscription struct {
	ID                uint   `gorm:"primaryKey"`
	UserID            uint   `gorm:"not null"`
	User              User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	City              string `gorm:"not null"`
	Frequency         string `gorm:"not null"`
	Status            string `gorm:"not null"`
	ConfirmationToken string `gorm:"not null"`
	UnsubscribeToken  string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
