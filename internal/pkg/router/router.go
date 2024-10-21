package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/slipneff/auto-yt/internal/utils/config"
)

type Handler struct {
	validator *validator.Validate
	cfg       *config.Config
}

func (h *Handler) InitRoutes() *gin.Engine {
	api := gin.New()

	return api
}

func NewRouter(
	cfg *config.Config,
) *Handler {
	return &Handler{
		validator: validator.New(),
		cfg:       cfg,
	}
}
