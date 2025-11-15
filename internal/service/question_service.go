package service

import (
	"context"
	"errors"

	"question-service/internal/domain"
	"question-service/internal/repository"

	"gorm.io/gorm"
)

type QuestionService struct {
	questions repository.QuestionRepository
}

func NewQuestionService(qRepo repository.QuestionRepository) *QuestionService {
	return &QuestionService{questions: qRepo}
}

// CreateQuestion создает новый вопрос.
func (s *QuestionService) CreateQuestion(ctx context.Context, text string) (*domain.Question, error) {
	q := &domain.Question{
		Text: text,
	}

	if err := s.questions.Create(ctx, q); err != nil {
		return nil, err
	}

	return q, nil
}

// ListQuestions возвращает список всех вопросов.
func (s *QuestionService) ListQuestions(ctx context.Context) ([]domain.Question, error) {
	return s.questions.GetAll(ctx)
}

// GetQuestionWithAnswers возвращает вопрос и все его ответы
func (s *QuestionService) GetQuestionWithAnswers(ctx context.Context, id int) (*domain.Question, error) {
	q, err := s.questions.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}

	return q, nil
}

// DeleteQuestion удаляет вопрос (каскад по FK удалит ответы).
func (s *QuestionService) DeleteQuestion(ctx context.Context, id int) error {
	err := s.questions.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return err
	}
	return nil
}
