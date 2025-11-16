package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"question-service/internal/domain"
	httptransport "question-service/internal/http"
	"question-service/internal/service"

	"github.com/stretchr/testify/require"
)

type mockAnswerRepo struct {
	nextID int
	answer []domain.Answer
}

func (f *mockAnswerRepo) Create(_ context.Context, q *domain.Answer) error {
	if f.nextID == 0 {
		f.nextID = 1
	}
	q.ID = f.nextID
	f.nextID++
	f.answer = append(f.answer, *q)
	return nil
}

func (f *mockAnswerRepo) GetByID(_ context.Context, id int) (*domain.Answer, error) {
	for _, q := range f.answer {
		if q.ID == id {
			qq := q
			return &qq, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *mockAnswerRepo) Delete(_ context.Context, id int) error {
	return nil
}

func (f *mockAnswerRepo) ListByQuestionID(_ context.Context, questionID int) ([]domain.Answer, error) {
	return nil, nil
}
func TestCreateAnswer_Success(t *testing.T) {
	aRepo := &mockAnswerRepo{}
	qRepo := &mockQuestionRepo{
		questions: []domain.Question{
			{ID: 1, Text: "What is GORM?"},
		},
	}
	svc := service.NewAnswerService(aRepo, qRepo)
	handler := httptransport.NewAnswerHandler(svc)

	body := []byte(`{"user_id":"user-123","text":"Answer text"}`)
	req := httptest.NewRequest(http.MethodPost, "/questions/1/answers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleCreateForQuestion(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var got domain.Answer
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)

	require.Equal(t, 1, got.QuestionID)
	require.Equal(t, "user-123", got.UserID)
	require.Equal(t, "Answer text", got.Text)
	require.NotZero(t, got.ID)
}

func TestCreateAnswer_QuestionNotFound(t *testing.T) {

	qRepo := &mockQuestionRepo{}
	aRepo := &mockAnswerRepo{}

	svc := service.NewAnswerService(aRepo, qRepo)
	handler := httptransport.NewAnswerHandler(svc)

	body := []byte(`{"user_id":"user-123","text":"Answer text"}`)
	req := httptest.NewRequest(http.MethodPost, "/questions/999/answers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleCreateForQuestion(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errResp struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	require.Equal(t, "question not found", errResp.Error)
}
