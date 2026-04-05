package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
)

func RequireAgent(userClient client.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx := auth.GetAuth(c)
		if authCtx == nil {
			c.Error(errors.UnauthorizedErr("not authenticated"))
			c.Abort()
			return
		}

		if authCtx.IdentityType != auth.IdentityEmployee || authCtx.EmployeeID == nil {
			c.Error(errors.ForbiddenErr("only employees can access this resource"))
			c.Abort()
			return
		}

		resp, err := userClient.GetEmployeeById(c.Request.Context(), uint64(*authCtx.EmployeeID))

		if err != nil {
			c.Error(errors.InternalErr(err))
			c.Abort()
			return
		}

		if !resp.IsAgent {
			c.Error(errors.ForbiddenErr("only agents can access this resource"))
			c.Abort()
			return
		}

		c.Next()
	}
}
