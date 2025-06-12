package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"url-shortener/config"
	"url-shortener/domains/urls/models/requests"
	"url-shortener/domains/urls/models/responses"
	"url-shortener/domains/urls/usecases"
	"url-shortener/shared"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UrlMappingHandler struct {
	UrlMappingUsecase *usecases.UrlMappingUsecase
	Validator         *validator.Validate
}

func NewUrlMappingHandler(urlMappingUsecase *usecases.UrlMappingUsecase) *UrlMappingHandler {
	return &UrlMappingHandler{
		UrlMappingUsecase: urlMappingUsecase,
		Validator:         validator.New(),
	}
}

func (h *UrlMappingHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var req requests.UrlMappingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendErrorResponse(w, shared.NewBadRequestError("Invalid request payload"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		shared.SendErrorResponse(w, shared.NewBadRequestError("Invalid request data"))
		return
	}

	mapping, err := h.UrlMappingUsecase.ShortenUrl(req.LongUrl, nil)
	if err != nil {
		log.Printf("Error in ShortenUrl usecase: %v", err)
		if respErr, ok := err.(shared.ErrorResponse); ok {
			shared.SendErrorResponse(w, respErr)
		} else {
			shared.SendErrorResponse(w, shared.NewInternalServerError("Unexpected error in shorten URL"))
		}
		return
	}

	baseURLStr := config.GetConfig().GetString("app.short_url")
	if baseURLStr == "" {
		shared.SendErrorResponse(w, shared.NewInternalServerError("Base URL configuration not found"))
		return
	}

	response := responses.UrlMappingResponse{
		ShortCode: mapping.ShortCode,
		LongUrl:   mapping.LongURL, // Make sure this contains the original long URL
		ExpiresAt: mapping.ExpiresAt,
		FullUrl:   baseURLStr + mapping.ShortCode,
	}

	shared.SendSuccessResponse(w, response, "Short URL created successfully")
}

func (h *UrlMappingHandler) GetLongUrlFromCode(c *gin.Context) {
	shortCode := c.Query("short_code")

	// Validasi shortCode
	if shortCode == "" {
		shared.SendErrorResponse(c.Writer, shared.NewBadRequestError("Short code is required"))
		return
	}

	mapping, err := h.UrlMappingUsecase.GetOriUrlByShortCode(shortCode)
	if err != nil {
		if respErr, ok := err.(shared.ErrorResponse); ok {
			shared.SendErrorResponse(c.Writer, respErr)
		} else {
			shared.SendErrorResponse(c.Writer, shared.NewInternalServerError("Unexpected error"))
		}
		return
	}

	baseURLStr := config.GetConfig().GetString("app.short_url")
	if baseURLStr == "" {
		shared.SendErrorResponse(c.Writer, shared.NewInternalServerError("Base URL configuration not found"))
		return
	}

	// Format respons JSON
	response := responses.UrlMappingResponse{
		ShortCode: mapping.ShortCode,
		LongUrl:   mapping.LongURL,
		ExpiresAt: mapping.ExpiresAt,
		FullUrl:   baseURLStr + mapping.ShortCode,
	}

	shared.SendSuccessResponse(c.Writer, response, "URL mapping retrieved successfully")
}

// Handler untuk redirect dengan clean URL pattern
func (h *UrlMappingHandler) RedirectOriUrl(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	// Validasi shortCode
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Get original URL from database
	mapping, err := h.UrlMappingUsecase.GetOriUrlByShortCode2(shortCode, ipAddress, userAgent)
	if err != nil {
		if respErr, ok := err.(shared.ErrorResponse); ok {
			switch respErr.GetStatus() {
			case 404:
				http.Error(w, "Short URL not found", http.StatusNotFound)
			case 400:
				http.Error(w, "Short URL has expired", http.StatusGone)
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, mapping.LongURL, http.StatusMovedPermanently)
}
