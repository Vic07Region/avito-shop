package mw

import (
	"context" //nolint:gci
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/Vic07Region/avito-shop/internal/storage" //nolint:gci
	"github.com/Vic07Region/avito-shop/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserStorage interface {
	GetUser4UserID(ctx context.Context, userID uuid.UUID) (*storage.Employee, error)
}

type Middleware struct {
	UserStorage
}

func New(userStorage UserStorage) *Middleware {
	return &Middleware{userStorage}
}

type User struct {
	UserID uuid.UUID
}

func (mw *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Unauthorized, Missing token"})
			c.Abort()
			return
		}
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid token type"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(strings.TrimPrefix(tokenString, "Bearer "))
		if err != nil {
			if errors.Is(err, utils.ErrorExpiredOrNotActive) {
				c.JSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "Invalid token"})
			c.Abort()
			return
		}

		_, err = mw.UserStorage.GetUser4UserID(context.Background(), claims.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "storage error"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
