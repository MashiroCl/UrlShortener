package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mashirocl/urlshortener/internal/model"
)

type URLService interface {
	CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error)
	GetURL(ctx context.Context, shortCode string) (string, error)
}

type URLHandler struct {
	urlService URLService
}

func NewUrlHandler(urlService URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// POST /api/url original_url, custom_code, duration -> short_code, duration
func (h *URLHandler) CreateURL(c echo.Context) error {
	fmt.Println("print start create url")
	// extract data
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	log.Println("data extracted")
	// validate data
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}
	log.Println("data verified")
	//transfer to short_code
	resp, err := h.urlService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

// GET /:code redirect code to original_url
func (h *URLHandler) RedirectURL(c echo.Context) error {
	// extract shortcode
	shortCode := c.Param("code")
	// search in database
	originalURL, err := h.urlService.GetURL(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, err.Error())
	}
	return c.Redirect(http.StatusPermanentRedirect, originalURL)
}
