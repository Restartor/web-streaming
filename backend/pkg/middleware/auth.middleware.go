package middleware

import (
	"backend/pkg/response"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// take token from header auth
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header tidak ada!")
			c.Abort()
			return
		}
		// format : bearer <token>, harus split bearer
		tokenString := strings.Split(authHeader, " ")

		if len(tokenString) != 2 {
			response.Error(c, http.StatusUnauthorized, "format token salah/error")
			c.Abort()
			return
		}

		// cek apakah sudah logout

		// parsing jwt.parse
		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
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
