package http

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"question-service/internal/logger"
	"strconv"
	"strings"

	"question-service/internal/service"
	"question-service/internal/transport"
)

type AnswerHandler struct {
	svc *service.AnswerService
	log *logger.Logger
}

func NewAnswerHandler(svc *service.AnswerService, log *logger.Logger) *AnswerHandler {
	return &AnswerHandler{svc: svc, log: log}
}

// HandleCreateForQuestion обрабатывает POST /questions/{id}/answers
func (h *AnswerHandler) HandleCreateForQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	path := strings.TrimPrefix(r.URL.Path, "/questions/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "answers" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		h.log.Warn("invalid question id",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("id_raw", parts[0]),
		)
		transport.WriteError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid json in answer creation",
			zap.Error(err),
			zap.String("path", r.URL.Path),
		)
		transport.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.TrimSpace(req.UserID) == "" || strings.TrimSpace(req.Text) == "" {
		h.log.Warn("missing required fields in answer creation",
			zap.String("user_id", req.UserID),
			zap.String("text", req.Text),
		)
		transport.WriteError(w, http.StatusBadRequest, "user_id and text are required")
		return
	}

	ans, err := h.svc.CreateAnswer(r.Context(), id, req.UserID, req.Text)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			h.log.Info("attempt to create answer for non-existing question",
				zap.Int("question_id", id),
				zap.String("user_id", req.UserID),
			)
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		h.log.Error("failed to create answer",
			zap.Error(err),
			zap.Int("question_id", id),
			zap.String("user_id", req.UserID),
		)
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
		h.log.Warn("invalid answer id",
			zap.String("path", r.URL.Path),
			zap.String("id_raw", idStr),
		)
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
			h.log.Info("answer not found", zap.Int("answer_id", id))
			transport.WriteError(w, http.StatusNotFound, "answer not found")
			return
		}

		h.log.Error("failed to get answer", zap.Error(err), zap.Int("answer_id", id))
		transport.WriteError(w, http.StatusInternalServerError, "failed to get answer")
		return
	}

	h.log.Info("answer fetched", zap.Int("answer_id", ans.ID))
	transport.WriteJSON(w, http.StatusOK, ans)
}

func (h *AnswerHandler) deleteAnswer(w http.ResponseWriter, r *http.Request, id int) {
	err := h.svc.DeleteAnswer(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			h.log.Info("attempt to delete non-existing answer", zap.Int("answer_id", id))
			transport.WriteError(w, http.StatusNotFound, "answer not found")
			return
		}

		h.log.Error("failed to delete answer", zap.Error(err), zap.Int("answer_id", id))
		transport.WriteError(w, http.StatusInternalServerError, "failed to delete answer")
		return
	}

	h.log.Info("answer deleted", zap.Int("answer_id", id))
	w.WriteHeader(http.StatusNoContent)
}
