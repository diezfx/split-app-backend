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

func newAPIHandler(projectService ProjectService) *ApiHandler {
	return &ApiHandler{projectService: projectService}
}

func InitAPI(_ config.Config, projectService ProjectService) *http.Server {
	mr := gin.New()
	mr.Use(gin.Recovery())

	r := mr.Group("/api/v1.0/")
	r.Use(middleware.HTTPLoggingMiddleware())

	apiHandler := newAPIHandler(projectService)

	r.GET("projects/:id", apiHandler.getProjectByIDHandler)
	r.GET("projects", apiHandler.getProjectsHandler)
	r.POST("projects", apiHandler.addProjectHandler)
	r.POST("projects/:id/transactions", apiHandler.addTransactionHandler)

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

func (api *ApiHandler) getProjectsHandler(ctx *gin.Context) {

	var queryParams GetProjectsQueryParams
	err := ctx.BindQuery(&queryParams)
	if err != nil {
		handleError(ctx, fmt.Errorf("getProjectHandler parse query params : %w: %w", errInvalidInput, err))
		return
	}

	proj, err := api.projectService.GetProjects(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, proj)
}

func (api *ApiHandler) getProjectByIDHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse id: %w: %w", errInvalidInput, err))
		return
	}

	proj, err := api.projectService.GetProjectByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, proj)
}

func (api *ApiHandler) addTransactionHandler(ctx *gin.Context) {

	var transaction AddTransaction
	if err := ctx.BindJSON(&transaction); err != nil {
		handleError(ctx, fmt.Errorf("parse add transaction body: %w: %w", errInvalidInput, err))
		return
	}

	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse id: %w: %w", errInvalidInput, err))
		return
	}

	svcTransaction, err := transaction.Validate()
	if err != nil {
		handleError(ctx, fmt.Errorf("validate transaction: %w: %w", errInvalidInput, err))
		return
	}

	err = api.projectService.AddTransaction(ctx, id, svcTransaction)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (api *ApiHandler) addProjectHandler(ctx *gin.Context) {
	var body AddProject
	err := ctx.BindJSON(&body)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse add project body: %w: %w", errInvalidInput, err))
		return
	}

	idParsed, err := uuid.Parse(body.ID)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse id in body: %w: %w", errInvalidInput, err))
		return
	}

	project := service.Project{ID: idParsed, Name: body.Name, Members: body.Members}

	proj, err := api.projectService.AddProject(ctx, project)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, proj)
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
		logger.Error(ctx, err).Msg("unexpected error occurred")
		ctx.String(http.StatusInternalServerError, "unexpected error")
	}
}
