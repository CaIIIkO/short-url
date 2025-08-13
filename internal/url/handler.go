package url

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type ServiceInterface interface {
	CreateLink(ctx context.Context, originalURL string) (*Link, error)
	RedirectAndLog(ctx context.Context, code string, ip string, userAgent string, referrer string) (*Link, error)
	Stats(ctx context.Context, code string) (*[]Click, *Link, error)
	LinksList(ctx context.Context) (*[]Link, error)
}

type Handler struct {
	service ServiceInterface
	baseURL string
}

func NewURLHandler(service ServiceInterface, baseURL string) *Handler {
	return &Handler{service: service, baseURL: baseURL}
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateLink(r.Context(), input.Original)
	if err != nil {
		http.Error(w, "failed to create link", http.StatusInternalServerError)
		return
	}

	resp := CreateLinkResponse{
		ID:        link.ID,
		UserID:    link.UserID,
		Original:  link.Original,
		ShortURL:  h.baseURL + link.ShortCode,
		ShortCode: link.ShortCode,
		CreatedAt: link.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/url/")
	if code == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	link, err := h.service.RedirectAndLog(r.Context(), code, r.RemoteAddr, r.UserAgent(), r.Referer())
	if err != nil {
		http.Error(w, "redirect error", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, link.Original, http.StatusFound)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/url/stats/")
	if code == "" {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	clicks, link, err := h.service.Stats(r.Context(), code)
	if err != nil {
		http.Error(w, "error get stats", http.StatusBadRequest)
		return
	}

	responses := StatsResponse{
		ID:         link.ID,
		Original:   link.Original,
		ShortURL:   h.baseURL + link.ShortCode,
		ShortCode:  link.ShortCode,
		CreatedAt:  link.CreatedAt,
		IsActive:   link.IsActive,
		TotalClick: len(*clicks),
		Clicks:     *clicks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

func (h *Handler) LinksList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	links, err := h.service.LinksList(r.Context())
	if err != nil {
		http.Error(w, "error get links list", http.StatusInternalServerError)
		return
	}

	responses := make([]LinksListResponse, 0, len(*links))
	for _, l := range *links {
		resp := LinksListResponse{
			ID:        l.ID,
			UserID:    l.UserID,
			Original:  l.Original,
			ShortCode: l.ShortCode,
			ShortURL:  h.baseURL + l.ShortCode,
			CreatedAt: l.CreatedAt,
		}
		responses = append(responses, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}
