package service

import (
	"context"
	"errors"
	"question-service/internal/domain"
	"question-service/internal/repository"

	"gorm.io/gorm"
)

type AnswerService struct {
	answers   repository.AnswerRepository
	questions repository.QuestionRepository
}

func NewAnswerService(aRepo repository.AnswerRepository, qRepo repository.QuestionRepository) *AnswerService {
	return &AnswerService{
		answers:   aRepo,
		questions: qRepo,
	}
}

// CreateAnswer добавляет ответ к вопросу.
func (s *AnswerService) CreateAnswer(ctx context.Context, questionID int, userID, text string) (*domain.Answer, error) {
	_, err := s.questions.GetByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}

		return nil, err
	}

	ans := &domain.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	if err := s.answers.Create(ctx, ans); err != nil {
		return nil, err
	}

	return ans, nil
}

// GetAnswer возвращает конкретный ответ по id.
func (s *AnswerService) GetAnswer(ctx context.Context, id int) (*domain.Answer, error) {
	a, err := s.answers.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnswerNotFound
		}
		return nil, err
	}
	return a, nil
}

// DeleteAnswer удаляет ответ.
func (s *AnswerService) DeleteAnswer(ctx context.Context, id int) error {
	err := s.answers.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnswerNotFound
		}
		return err
	}
	return nil
}
