package service_test

import (
	"context"
	"testing"

	"question-service/internal/domain"
	"question-service/internal/service"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockQuestionRepo struct {
	created []*domain.Question
	getByID func(id int) (*domain.Question, error)
}

func (m *mockQuestionRepo) Create(_ context.Context, q *domain.Question) error {
	m.created = append(m.created, q)
	// эмулируем, что БД проставила ID
	if q.ID == 0 {
		q.ID = 1
	}
	return nil
}

func (m *mockQuestionRepo) GetAll(_ context.Context) ([]domain.Question, error) {
	var out []domain.Question
	for _, q := range m.created {
		out = append(out, *q)
	}
	return out, nil
}

func (m *mockQuestionRepo) GetByID(_ context.Context, id int) (*domain.Question, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockQuestionRepo) Delete(_ context.Context, id int) error {
	return nil
}

func TestQuestionService_CreateQuestion(t *testing.T) {
	repo := &mockQuestionRepo{}
	svc := service.NewQuestionService(repo)

	ctx := context.Background()
	q, err := svc.CreateQuestion(ctx, "What is GORM?")
	require.NoError(t, err)
	require.NotNil(t, q)
	require.Equal(t, "What is GORM?", q.Text)
	require.NotZero(t, q.ID)
}
func TestQuestionService_GetQuestionWithAnswers_NotFound(t *testing.T) {
	repo := &mockQuestionRepo{
		getByID: func(id int) (*domain.Question, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewQuestionService(repo)

	ctx := context.Background()
	q, err := svc.GetQuestionWithAnswers(ctx, 123)

	require.Nil(t, q)
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrQuestionNotFound)
}
