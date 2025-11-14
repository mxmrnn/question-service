package domain

import "time"

type Answer struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID int       `gorm:"not null;index"           json:"question_id"`
	UserID     string    `gorm:"type:varchar(64);not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null"        json:"text"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"   json:"created_at"`
}
