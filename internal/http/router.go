package http

import (
	"encoding/json"
	"net/http"

	"question-service/internal/service"
)

func NewRouter(qSvc *service.QuestionService, aSvc *service.AnswerService) *http.ServeMux {
	mux := http.NewServeMux()

	// healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	qh := NewQuestionHandler(qSvc)
	ah := NewAnswerHandler(aSvc)

	// /questions (GET, POST)
	mux.HandleFunc("/questions", qh.HandleQuestions)

	mux.HandleFunc("/questions/", func(w http.ResponseWriter, r *http.Request) {
		// /questions (GET, POST)
		if r.URL.Path == "/questions/" {
			qh.HandleQuestions(w, r)
			return
		}

		// /questions/{id}/answers (POST)
		if len(r.URL.Path) > len("/questions/") && hasSuffix(r.URL.Path, "/answers") {
			ah.HandleCreateForQuestion(w, r)
			return
		}

		// /questions/{id} (GET, DELETE)
		qh.HandleQuestionByID(w, r)
	})

	// /answers/{id} (GET, DELETE)
	mux.HandleFunc("/answers/", ah.HandleAnswerByID)

	return mux
}

// вместо strings.HasSuffix
func hasSuffix(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
