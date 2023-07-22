package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/diezfx/split-app-backend/internal/config"
	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/diezfx/split-app-backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApiHandler struct {
	projectService ProjectService
}

func newApiHandler(projectService ProjectService) *ApiHandler {
	return &ApiHandler{projectService: projectService}
}

func InitApi(cfg config.Config, projectService ProjectService) *http.Server {
	mr := gin.New()
	mr.Use(gin.Recovery())

	r := mr.Group("/api/v1.0/")
	r.Use(middleware.HTTPLoggingMiddleware())

	apiHandler := newApiHandler(projectService)

	r.GET("project/:id", apiHandler.getProjectHandler)

	mr.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "OPTION"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	return &http.Server{
		Handler: mr,
		Addr:    "localhost:5002",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (api *ApiHandler) getProjectHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse uuid: %w: %w", errInvalidInput, err))
		return
	}

	proj, err := api.projectService.GetProject(ctx, uuid)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, proj)
}

func handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, errInvalidInput):
		logger.Info(ctx).Err(err).Msg("request failed with invalid input")
		ctx.String(http.StatusBadRequest, "invalid input")
	case errors.Is(err, service.ErrProjectNotFound):
		logger.Info(ctx).Err(err).Msg("not found")
		ctx.String(http.StatusNotFound, "not found")
	default:
		logger.Error(ctx, err).Msg("unexpected error occured")
		ctx.String(http.StatusInternalServerError, "unexpected error")
	}
}
