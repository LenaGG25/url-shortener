package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"url-shortener/client/internal/service"
)

type Handler struct {
	s *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) CreateURL(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request createShortURLRequest
	if err := json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortUrl, err := h.s.CreateShortUrl(req.Context(), request.OriginalURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &createShortURLResponse{ShortURL: shortUrl}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetOriginalByShortURL(w http.ResponseWriter, req *http.Request) {
	shortURL, ok := mux.Vars(req)["short_url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, err := h.s.GetOriginalUrl(req.Context(), shortURL)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, originalURL, http.StatusMovedPermanently)
}

func (h *Handler) GetStatsByShortURL(w http.ResponseWriter, req *http.Request) {
	shortURL, ok := mux.Vars(req)["short_url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestNumber, err := h.s.GetStatistics(req.Context(), shortURL)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &getStatsByShortURLResponse{RequestNumber: requestNumber}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respJSON); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
