package access

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// DirectorReadOnly запрещает директору проектов (dp) изменять данные модуля:
// любые методы, кроме GET, отдаются 403. Директор может менять только приоритет
// проектов — это ограничение реализовано на уровне сервиса проекта.
func DirectorReadOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := helpers.GetUser(c)
		if err != nil {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		if user.Role == userdomain.ProjectDirector && c.Request.Method != http.MethodGet {
			response.Forbidden(c, "director (dp) is read-only")
			c.Abort()
			return
		}

		c.Next()
	}
}
