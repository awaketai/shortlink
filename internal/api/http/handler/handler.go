package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"shortlink/internal/shortener"
	"strings"
)

type LinkAPI struct {
	service *shortener.Service
	logger  *log.Logger
}

func NewLinkAPI(service *shortener.Service, l *log.Logger) *LinkAPI {
	if l == nil {
		logger := log.New(os.Stdout, "[LinkAPI] ", log.LstdFlags|log.Lshortfile)
		return &LinkAPI{
			service: service,
			logger:  logger,
		}
	}
	return &LinkAPI{
		service: service,
		logger:  l,
	}
}

type CreateShortLinkRequest struct {
	LongURL string `json:"long_url"`
}

type CreateShortLinkResponse struct {
	ShortCode string `json:"short_code"`
}

// CreateLink 创建短链接 POST请求
func (l *LinkAPI) CreateLink(w http.ResponseWriter, r *http.Request) {
	// 获取上下文
	ctx := r.Context()
	if r.Method != http.MethodPost {
		l.logger.Printf("WARN: Invalid method for create link: %s from %s\n", r.Method, r.RemoteAddr)
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.logger.Printf("ERROR: Failed to decode request body: %v\n", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if strings.TrimSpace(req.LongURL) == "" {
		l.logger.Printf("ERROR: Long URL is empty")
		http.Error(w, "Long URL is empty", http.StatusBadRequest)
		return
	}
	l.logger.Printf("INFO: Received request to create short link from %s, LongURL: %s\n", r.RemoteAddr, req.LongURL)

	shortCode, err := l.service.CreateShortLink(ctx, req.LongURL)
	if err != nil {
		l.logger.Printf("ERROR: Failed to create short link: %v\n", err)
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
		return
	}
	resp := CreateShortLinkResponse{
		ShortCode: shortCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		l.logger.Printf("ERROR: Failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	l.logger.Printf("INFO: Created short link %s for long URL %s\n", shortCode, req.LongURL)
}

func (l *LinkAPI) RedirectLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	l.logger.Printf("INFO: Received request to redirect short link from %s. ShortCode: %s, Path: %s\n", r.RemoteAddr, shortCode, r.URL.Path)

	if r.Method != http.MethodGet {
		l.logger.Printf("ERROR: Only GET method is allowed")
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// 基础路径检查，避免匹配到 /api/links,/healthz等
	if shortCode == "" || shortCode == "api/links" || shortCode == "healthz" {
		l.logger.Printf("INFO: Path is not a shortcode, treating as not found. Path: %s, from %s\n", r.URL.Path, r.RemoteAddr)
		http.NotFound(w, r)
		return
	}
	longURL, err := l.service.GetAndTrackLongURL(ctx, shortCode)
	l.logger.Printf("WARN: Service failed to get long URL for redirect from %s. ShortCode: %s, Error: %v\n", r.RemoteAddr, shortCode, err)
	// 具体错误处理
	if err != nil {
		l.logger.Printf("ERROR: Failed to get long URL for redirect: %v\n", err)
		http.Error(w, "Failed to get long URL for redirect", http.StatusInternalServerError)
		return
	}
	l.logger.Printf("INFO: Redirecting %s from %s to %s\n", shortCode, r.RemoteAddr, longURL)
	http.Redirect(w, r, longURL, http.StatusFound)
}
