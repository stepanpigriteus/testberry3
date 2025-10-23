package httpsh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"threeFour/domain"

	"github.com/rs/zerolog"
)

type HandlersImpl struct {
	serv   domain.Service
	logger zerolog.Logger
}

func NewHandlers(ctx context.Context, serv domain.Service, logger zerolog.Logger) *HandlersImpl {
	return &HandlersImpl{
		serv:   serv,
		logger: logger,
	}
}

func (h *HandlersImpl) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Info().Msg("Upload handler called")

	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse form file")
		writeError(w, http.StatusBadRequest, "incorrect file upload")
		return
	}
	defer file.Close()

	h.logger.Info().Str("filename", header.Filename).Msg("file received")

	data, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to read file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int("bytes", len(data)).Msg("file read complete")

	contentType := header.Header.Get("Content-Type")
	h.logger.Info().Str("contentType", contentType).Msg("uploading to service")

	id, err := h.serv.Upload(ctx, domain.ImageData{
		Bytes:       data,
		ContentType: contentType,
		Filename:    header.Filename,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("service upload failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info().Msg("upload completed successfully: id " + id)
	writeJSON(w, http.StatusCreated, "Image upload succesfully")
}

func (h *HandlersImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Info().Msg("Upload handler called")
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(path) != 2 {
		writeError(w, http.StatusBadRequest, "incorrect url path")
		return
	}
	id := path[len(path)-1]
	data, err := h.serv.Get(ctx, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get image")
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}
	contentType := http.DetectContentType(data)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (h *HandlersImpl) Delete(w http.ResponseWriter, r *http.Request) {

}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
