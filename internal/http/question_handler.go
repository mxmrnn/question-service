package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"question-service/internal/logger"
	"question-service/internal/service"
	"question-service/internal/transport"
)

type QuestionHandler struct {
	svc *service.QuestionService
	log *logger.Logger
}

func NewQuestionHandler(svc *service.QuestionService, log *logger.Logger) *QuestionHandler {
	return &QuestionHandler{svc: svc, log: log}
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
		h.log.Warn("invalid question id",
			zap.String("path", r.URL.Path),
			zap.String("id_raw", idStr),
		)
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
	defer r.Body.Close()

	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid json in create question",
			zap.Error(err),
			zap.String("path", r.URL.Path),
		)
		transport.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		h.log.Warn("empty text in create question",
			zap.String("path", r.URL.Path),
		)
		transport.WriteError(w, http.StatusBadRequest, "text is required")
		return
	}

	q, err := h.svc.CreateQuestion(r.Context(), req.Text)
	if err != nil {
		h.log.Error("failed to create question",
			zap.Error(err),
			zap.String("text", req.Text),
		)
		transport.WriteError(w, http.StatusInternalServerError, "failed to create question")
		return
	}

	h.log.Info("question created",
		zap.Int("question_id", q.ID),
	)

	transport.WriteJSON(w, http.StatusCreated, q)
}

func (h *QuestionHandler) listQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.svc.ListQuestions(r.Context())
	if err != nil {
		h.log.Error("failed to list questions",
			zap.Error(err),
		)
		transport.WriteError(w, http.StatusInternalServerError, "failed to list questions")
		return
	}

	h.log.Info("questions listed",
		zap.Int("count", len(questions)),
	)
	transport.WriteJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) getQuestion(w http.ResponseWriter, r *http.Request, id int) {
	q, err := h.svc.GetQuestionWithAnswers(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			h.log.Info("question not found",
				zap.Int("question_id", id),
			)
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		h.log.Error("failed to get question",
			zap.Error(err),
			zap.Int("question_id", id),
		)
		transport.WriteError(w, http.StatusInternalServerError, "failed to get question")
		return
	}

	h.log.Info("question fetched",
		zap.Int("question_id", q.ID),
	)

	transport.WriteJSON(w, http.StatusOK, q)
}

func (h *QuestionHandler) deleteQuestion(w http.ResponseWriter, r *http.Request, id int) {
	err := h.svc.DeleteQuestion(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			h.log.Info("attempt to delete non-existing question",
				zap.Int("question_id", id),
			)
			transport.WriteError(w, http.StatusNotFound, "question not found")
			return
		}
		h.log.Error("failed to delete question",
			zap.Error(err),
			zap.Int("question_id", id),
		)
		transport.WriteError(w, http.StatusInternalServerError, "failed to delete question")
		return
	}
	h.log.Info("question deleted",
		zap.Int("question_id", id),
	)
	w.WriteHeader(http.StatusNoContent)
}
