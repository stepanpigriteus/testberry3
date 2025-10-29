package httpsh

import (
	"archive/zip"
	"bytes"
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
	writeJSON(w, http.StatusCreated, fmt.Sprintf("Image upload successfully: id %s", id))

}

func (h *HandlersImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := strings.TrimPrefix(r.URL.Path, "/get/")
	path := strings.Split(strings.Trim(id, "/"), "/")
	if len(path) != 2 {
		writeError(w, http.StatusBadRequest, "incorrect url path")
		return
	}
	final := path[1]
	fmt.Println(final)
	files, err := h.serv.Get(ctx, final)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get images")
		http.Error(w, "failed to get images", http.StatusInternalServerError)
		return
	}
	if files == nil {
		http.Error(w, "image still processing", http.StatusAccepted)
		return
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	names := []string{"watermark.jpg", "resized.jpg", "thumbnail.jpg"}
	for i, data := range files {
		f, err := zipWriter.Create(names[i])
		if err != nil {
			h.logger.Err(err).Msg("failed to create zip entry")
			continue
		}
		if _, err := f.Write(data); err != nil {
			h.logger.Err(err).Msg("failed to write to zip entry")
			continue
		}
	}

	if err := zipWriter.Close(); err != nil {
		h.logger.Err(err).Msg("failed to close zip writer")
		http.Error(w, "failed to create zip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_images.zip\"", final))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (h *HandlersImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	h.logger.Info().Msg("Delete handler called")
	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(path) != 2 {
		writeError(w, http.StatusBadRequest, "incorrect url path")
		return
	}
	id := path[len(path)-1]
	err := h.serv.Delete(ctx, id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to delete image")
		writeError(w, http.StatusBadRequest, "failed to delete image")
		return
	}
	h.logger.Info().Msg("upload deleted successfully: id " + id)
	writeJSON(w, http.StatusOK, "Image deleted succesfully")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
