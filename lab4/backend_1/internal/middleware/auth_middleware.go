package middleware

import (
	logging "backend/internal"
	"backend/internal/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logging.Logger.Warn("Отсутствует токен аутентификации")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен аутентификации"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := services.ValidateToken(tokenString)
		if err != nil {
			logging.Logger.Warn("Неверный токен аутентификации: ", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logging.Logger.Warn("Не удалось извлечь данные из токена")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Не удалось извлечь данные из токена"})
			return
		}

		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		c.Set("userID", userID)
		c.Set("role", role)
		logging.Logger.Info("Аутентификация успешна, пользователь ID: ", userID)

		c.Next()
	}
}
