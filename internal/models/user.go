package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FirebaseUID  string    `json:"firebase_uid" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"not null"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
