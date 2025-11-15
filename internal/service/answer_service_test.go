package service_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"question-service/internal/service"

	"question-service/internal/domain"
	"testing"
)

type mockAnswerRepo struct {
	created []*domain.Answer
}

func (m *mockAnswerRepo) Create(_ context.Context, a *domain.Answer) error {
	if a.ID == 0 {
		a.ID = 1
	}
	m.created = append(m.created, a)
	return nil
}

func (m *mockAnswerRepo) GetByID(_ context.Context, id int) (*domain.Answer, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAnswerRepo) Delete(_ context.Context, id int) error { return nil }

func (m *mockAnswerRepo) ListByQuestionID(_ context.Context, qid int) ([]domain.Answer, error) {
	return nil, nil
}

func TestAnswerService_CreateAnswer_QuestionNotFound(t *testing.T) {

	qRepo := &mockQuestionRepo{
		getByID: func(id int) (*domain.Question, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	aRepo := &mockAnswerRepo{}
	svc := service.NewAnswerService(aRepo, qRepo)

	ctx := context.Background()
	ans, err := svc.CreateAnswer(ctx, 999, "u123", "hi")

	require.Nil(t, ans)
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrQuestionNotFound)
}
