package domain

import "time"

type Question struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Text      string    `gorm:"type:text;not null"       json:"text"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"  json:"created_at"`
	Answers   []Answer  `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}
