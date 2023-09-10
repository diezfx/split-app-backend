package auth

import (
	"net/http"
	"strings"

	"github.com/diezfx/split-app-backend/internal/contextutil"
	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

const AuthorizationHeader = "Authorization"

const BearerPrefix = "Bearer "

func AuthMiddleware(cfg Config) gin.HandlerFunc {
	authValidator := New(cfg)
	return func(ctx *gin.Context) {
		authToken := ctx.Request.Header.Get(AuthorizationHeader)
		if authToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authTokenStripped, _ := strings.CutPrefix(authToken, BearerPrefix)
		token, err := authValidator.Validate(authTokenStripped)
		if err != nil {
			logger.Debug(ctx).Err(err).Msg("validation of bearer token failed")
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		sub := token.Subject()
		requestCtx := contextutil.AddUserIDToCtx(ctx.Request.Context(), sub)
		ctx.Request = ctx.Request.WithContext(requestCtx)
		ctx.Next()
	}
}
