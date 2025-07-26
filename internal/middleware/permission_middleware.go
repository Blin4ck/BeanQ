package middleware

import (
	"coffe/internal/user/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PermissionMiddleware(permissionUC usecase.PermissionUsecase, required string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleIDRaw, exists := ctx.Get("role_id") // допустим, роль уже есть в контексте
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
			return
		}

		roleID, ok := roleIDRaw.(uuid.UUID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid role"})
			return
		}

		perms, err := permissionUC.GetPermissionsByRole(ctx.Request.Context(), roleID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get permissions"})
			return
		}

		for _, p := range perms {
			if strings.EqualFold(p.Code(), required) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
