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

type AnswerHandler struct {
	svc *service.AnswerService
}

func NewAnswerHandler(svc *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{svc: svc}
}

// HandleCreateForQuestion обрабатывает POST /questions/{id}/answers
func (h *AnswerHandler) HandleCreateForQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/questions/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "answers" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		transport.WriteError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Text) == "" {
		transport.WriteError(w, http.StatusBadRequest, "user_id and text are required")
		return
	}

	ans, err := h.svc.CreateAnswer(r.Context(), id, req.UserID, req.Text)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		transport.WriteError(w, http.StatusInternalServerError, "failed to create answer")
		return
	}

	transport.WriteJSON(w, http.StatusCreated, ans)
}

// HandleAnswerByID обрабатывает GET /answers/{id} и DELETE /answers/{id}
func (h *AnswerHandler) HandleAnswerByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/answers/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		transport.WriteError(w, http.StatusBadRequest, "invalid answer id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAnswer(w, r, id)
	case http.MethodDelete:
		h.deleteAnswer(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *AnswerHandler) getAnswer(w http.ResponseWriter, r *http.Request, id int) {
	ans, err := h.svc.GetAnswer(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			transport.WriteError(w, http.StatusNotFound, "answer not found")
			return
		}
		transport.WriteError(w, http.StatusInternalServerError, "failed to get answer")
		return
	}

	transport.WriteJSON(w, http.StatusOK, ans)
}

func (h *AnswerHandler) deleteAnswer(w http.ResponseWriter, r *http.Request, id int) {
	err := h.svc.DeleteAnswer(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			transport.WriteError(w, http.StatusNotFound, "answer not found")
			return
		}
		transport.WriteError(w, http.StatusInternalServerError, "failed to delete answer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
