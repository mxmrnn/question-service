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

type mockQuestionRepo struct {
	nextID    int
	questions []domain.Question
}

func (f *mockQuestionRepo) Create(_ context.Context, q *domain.Question) error {
	if f.nextID == 0 {
		f.nextID = 1
	}
	q.ID = f.nextID
	f.nextID++
	f.questions = append(f.questions, *q)
	return nil
}

func (f *mockQuestionRepo) GetAll(_ context.Context) ([]domain.Question, error) {
	return f.questions, nil
}

func (f *mockQuestionRepo) GetByID(_ context.Context, id int) (*domain.Question, error) {
	for _, q := range f.questions {
		if q.ID == id {
			qq := q
			return &qq, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *mockQuestionRepo) Delete(_ context.Context, id int) error {
	return nil
}

func TestCreateQuestion_Success(t *testing.T) {
	repo := &mockQuestionRepo{}
	svc := service.NewQuestionService(repo)
	handler := httptransport.NewQuestionHandler(svc)

	body := []byte(`{"text":"What is GORM?"}`)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleQuestions(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var got domain.Question
	err := json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err)
	require.Equal(t, "What is GORM?", got.Text)
	require.NotZero(t, got.ID)
}
func TestCreateQuestion_InvalidJSON(t *testing.T) {
	repo := &mockQuestionRepo{}
	svc := service.NewQuestionService(repo)
	handler := httptransport.NewQuestionHandler(svc)

	body := []byte(`{"text":`)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleQuestions(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errResp struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	require.Equal(t, "invalid json", errResp.Error)
}
func TestCreateQuestion_EmptyText(t *testing.T) {
	repo := &mockQuestionRepo{}
	svc := service.NewQuestionService(repo)
	handler := httptransport.NewQuestionHandler(svc)

	body := []byte(`{"text":""}`)
	req := httptest.NewRequest(http.MethodPost, "/questions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.HandleQuestions(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errResp struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	require.Equal(t, "text is required", errResp.Error)
}
