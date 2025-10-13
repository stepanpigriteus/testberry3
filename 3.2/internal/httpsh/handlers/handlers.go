package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"treeTwo/domain"
	"treeTwo/pkg"

	"github.com/rs/zerolog"
)

type ShortHandlers struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *ShortHandlers {
	return &ShortHandlers{
		serv:   serv,
		logger: logger,
	}
}

func (h *ShortHandlers) CreateShorten(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("CreateShorten handler called")
	ctx := context.Background()
	var link domain.Link

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	if !strings.HasPrefix(link.Line, "http://") && !strings.HasPrefix(link.Line, "https://") {
		writeError(w, http.StatusBadRequest, "ссылка должна начинаться с http:// или https://")
		return
	}
	if err := h.serv.CreateShorten(ctx, link.Line); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "уже существует") {
			status = http.StatusConflict
		}
		writeError(w, status, fmt.Sprintf("ошибка создания ссылки: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *ShortHandlers) GetShorten(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("GetShorten handler called")
	ctx := context.Background()

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // убираем лишние "/"
	if len(path) < 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	code := path[1]
	h.logger.Info().Msgf("short code: %s (len=%d)", code, len(code))

	if len(code) != 6 {
		writeError(w, http.StatusBadRequest, "link length must equal 6")
		return
	}

	link, err := h.serv.GetShorten(ctx, code)
	if err != nil {
		writeError(w, http.StatusBadRequest, "не получилось вытащить")
		return
	}

	visit := domain.Visit{
		ID:          link.ID,
		ShortLinkID: link.ID,
		VisitedAt:   time.Now(),
		UserAgent:   r.UserAgent(),
		IPAddress:   pkg.GetIP(r),
		DeviceType:  pkg.DetectDevice(r.UserAgent()),
	}

	if err := h.serv.SaveVisit(ctx, visit); err != nil {
		h.logger.Err(err).Msg("не получилось сохранить визит")
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func (h *ShortHandlers) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	path := strings.Split(r.URL.Path, "/")
	if len(path[2]) != 6 {
		writeError(w, http.StatusBadRequest, "link length must equal 6")
		return
	}
	group := r.URL.Query().Get("group")
	if group != "day" && group != "month" && group != "detailed" {
		h.logger.Error().Msg("incorrect group")
		writeError(w, http.StatusBadRequest, "incorrect query parameter")
	}
	stat, err := h.serv.GetAnalytics(ctx, path[2], group)
	if err != nil {
		// errors.Is() позже
		writeError(w, http.StatusBadRequest, "не получилось вытащить")
		return
	}
	writeJSON(w, http.StatusOK, stat)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
