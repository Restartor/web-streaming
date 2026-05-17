package middleware

import (
	"backend/config"
	"backend/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(cfg config.AppConfig) gin.HandlerFunc {

	return func(c *gin.Context) {
		// take token from header auth
		authHeader, err := c.Cookie("access_token")
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Authorization header tidak ada!")
			c.Abort()
			return
		}

		// cek apakah sudah logout

		// parsing jwt.parse
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "token not valid")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(userID))
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
		}

		c.Next()

	}

}
