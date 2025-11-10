package httpsh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"threeSixth/domain"
	"time"

	"github.com/rs/zerolog"
)

type Handlers struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *Handlers {
	return &Handlers{
		serv:   serv,
		logger: logger,
	}
}

func (h *Handlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Create Item called")
	ctx := r.Context()
	var item domain.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "не удалось прочитать JSON: "+err.Error())
		return
	}
	defer r.Body.Close()

	data, err := h.serv.CreateItem(ctx, item)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed create")
	}

	writeJSON(w, http.StatusCreated, data)
}

func (h *Handlers) GetItems(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("GetItem handler called")
	ctx := r.Context()

	filter := domain.Filter{
		From:     r.URL.Query().Get("from"),
		To:       r.URL.Query().Get("to"),
		Type:     r.URL.Query().Get("type"),
		Category: r.URL.Query().Get("category"),
		SortBy:   r.URL.Query().Get("sort_by"),
		Order:    r.URL.Query().Get("order"),
	}
	fmt.Println(filter)
	items, err := h.serv.GetItems(ctx, filter)
	if err != nil {
		writeError(w, http.StatusBadRequest, "incorrect request parameters")
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handlers) GetItem(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("GetItem handler called")

}

func (h *Handlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/items/")
	if path == "" {
		writeError(w, http.StatusBadRequest, "empty id")
	}
	var input domain.Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	item, err := h.serv.UpdateItem(ctx, path, input)
	if err != nil {
		h.logger.Err(err).Msg("update item failed")
		writeError(w, http.StatusInternalServerError, err.Error())
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Handlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("DeleteItem called")
	ctx := r.Context()

	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	err := h.serv.DeleteItem(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handlers) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("GetAnalytics handler called")
	ctx := r.Context()

	filter := domain.AnalyticsFilter{
		From:    r.URL.Query().Get("from"),
		To:      r.URL.Query().Get("to"),
		Type:    r.URL.Query().Get("type"),
		GroupBy: r.URL.Query().Get("group_by"),
	}

	if err := validateAnalyticsFilter(filter); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.serv.GetAnalytics(ctx, filter)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get analytics")
		writeError(w, http.StatusInternalServerError, "failed to get analytics")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func validateAnalyticsFilter(filter domain.AnalyticsFilter) error {

	if filter.From != "" {
		if _, err := time.Parse("2006-01-02", filter.From); err != nil {
			return fmt.Errorf("invalid 'from' date format, expected YYYY-MM-DD")
		}
	}

	if filter.To != "" {
		if _, err := time.Parse("2006-01-02", filter.To); err != nil {
			return fmt.Errorf("invalid 'to' date format, expected YYYY-MM-DD")
		}
	}

	if filter.Type != "" {
		validTypes := map[string]bool{"expense": true, "income": true}
		if !validTypes[filter.Type] {
			return fmt.Errorf("invalid type, must be 'expense' or 'income'")
		}
	}
	if filter.GroupBy != "" {
		validGroupBy := map[string]bool{"category": true, "date": true, "type": true}
		if !validGroupBy[filter.GroupBy] {
			return fmt.Errorf("invalid group_by, must be 'category', 'date', or 'type'")
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
