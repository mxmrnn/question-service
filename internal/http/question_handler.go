package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"question-service/internal/service"
	"question-service/internal/transport"
)

type QuestionHandler struct {
	svc *service.QuestionService
}

func NewQuestionHandler(svc *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{svc: svc}
}

// HandleQuestions обрабатывает /questions (GET, POST)
func (h *QuestionHandler) HandleQuestions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listQuestions(w, r)
	case http.MethodPost:
		h.createQuestion(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleQuestionByID обрабатывает /questions/{id} (GET, DELETE)
func (h *QuestionHandler) HandleQuestionByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/questions/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	if strings.Contains(idStr, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		transport.WriteError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getQuestion(w, r, id)
	case http.MethodDelete:
		h.deleteQuestion(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type createQuestionRequest struct {
	Text string `json:"text"`
}

func (h *QuestionHandler) createQuestion(w http.ResponseWriter, r *http.Request) {
	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		transport.WriteError(w, http.StatusBadRequest, "text is required")
		return
	}

	q, err := h.svc.CreateQuestion(r.Context(), req.Text)
	if err != nil {
		transport.WriteError(w, http.StatusInternalServerError, "failed to create question")
		return
	}

	transport.WriteJSON(w, http.StatusCreated, q)
}

func (h *QuestionHandler) listQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.svc.ListQuestions(r.Context())
	if err != nil {
		transport.WriteError(w, http.StatusInternalServerError, "failed to list questions")
		return
	}

	transport.WriteJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) getQuestion(w http.ResponseWriter, r *http.Request, id int) {
	q, err := h.svc.GetQuestionWithAnswers(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		transport.WriteError(w, http.StatusInternalServerError, "failed to get question")
		return
	}

	transport.WriteJSON(w, http.StatusOK, q)
}

func (h *QuestionHandler) deleteQuestion(w http.ResponseWriter, r *http.Request, id int) {
	err := h.svc.DeleteQuestion(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		transport.WriteError(w, http.StatusInternalServerError, "failed to delete question")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
