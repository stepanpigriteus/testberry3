package httpsh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"treethree/domain"

	"github.com/rs/zerolog"
)

type CommentsHandlers struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *CommentsHandlers {
	return &CommentsHandlers{
		serv:   serv,
		logger: logger,
	}
}

func (h *CommentsHandlers) CreateComments(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Info().Msg("Create handler called >>>")
	var comment domain.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	if comment.Text == "" && comment.Author == "" && comment.ParentID == nil {
		writeError(w, http.StatusBadRequest, "emty fields in json")
		return
	}

	err := h.serv.CreateComments(ctx, comment)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("ошибка создания ссылки: %v", err))
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}
func (h *CommentsHandlers) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Info().Msg("GetComments handler called >>>")
	query := r.URL.Query().Get("parent")
	id, err := strconv.Atoi(query)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("incorrect query parameter,  error : %v", err))
		return
	}
	comments, err := h.serv.GetComments(ctx, id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "не получилось вытащить")
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *CommentsHandlers) Delete(w http.ResponseWriter, r *http.Request) {

}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
