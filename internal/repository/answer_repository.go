package repository

import (
	"context"

	"gorm.io/gorm"
	"question-service/internal/domain"
)

type AnswerRepository interface {
	Create(ctx context.Context, a *domain.Answer) error
	GetByID(ctx context.Context, id int) (*domain.Answer, error)
	Delete(ctx context.Context, id int) error
	ListByQuestionID(ctx context.Context, questionID int) ([]domain.Answer, error)
}

type GormAnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *GormAnswerRepository {
	return &GormAnswerRepository{db: db}
}

func (r *GormAnswerRepository) Create(ctx context.Context, a *domain.Answer) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *GormAnswerRepository) GetByID(ctx context.Context, id int) (*domain.Answer, error) {
	var ans domain.Answer
	err := r.db.WithContext(ctx).First(&ans, id).Error
	if err != nil {
		return nil, err
	}
	return &ans, nil
}

func (r *GormAnswerRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Answer{}, id).Error
}

func (r *GormAnswerRepository) ListByQuestionID(ctx context.Context, questionID int) ([]domain.Answer, error) {
	var answers []domain.Answer
	err := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&answers).Error
	return answers, err
}
