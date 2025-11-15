package repository

import (
	"context"

	"gorm.io/gorm"
	"question-service/internal/domain"
)

type QuestionRepository interface {
	Create(ctx context.Context, q *domain.Question) error
	GetAll(ctx context.Context) ([]domain.Question, error)
	GetByID(ctx context.Context, id int) (*domain.Question, error)
	Delete(ctx context.Context, id int) error
}

type GormQuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *GormQuestionRepository {
	return &GormQuestionRepository{db: db}
}

func (r *GormQuestionRepository) Create(ctx context.Context, q *domain.Question) error {
	return r.db.WithContext(ctx).Create(q).Error
}

func (r *GormQuestionRepository) GetAll(ctx context.Context) ([]domain.Question, error) {
	var questions []domain.Question
	err := r.db.WithContext(ctx).Find(&questions).Error
	return questions, err
}

func (r *GormQuestionRepository) GetByID(ctx context.Context, id int) (*domain.Question, error) {
	var q domain.Question
	err := r.db.WithContext(ctx).
		Preload("Answers").
		First(&q, id).Error
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *GormQuestionRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&domain.Question{}, id).Error
}
