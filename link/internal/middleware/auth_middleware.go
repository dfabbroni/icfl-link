package middleware

import (
	"strings"
	"time"

	"link/internal/models"
	"link/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CombinedAuthMiddleware(db *gorm.DB, secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.NewUnauthorizedError("Missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.NewUnauthorizedError("Invalid authorization header format")
			}

			token := parts[1]

			// Try Node authentication first
			var node models.Node
			if err := db.Where("public_key = ?", token).First(&node).Error; err == nil {
				if !node.Approved {
					return utils.NewUnauthorizedError("Node not approved")
				}

				c.Set("node", node)
				return next(c)
			}

			// If Node authentication fails, try JWT authentication
			jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, utils.NewUnauthorizedError("Invalid token signing method")
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return utils.NewUnauthorizedError("Invalid or expired token")
			}

			if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					return utils.NewUnauthorizedError("Token has expired")
				}

				c.Set("user_id", claims["user_id"])
				return next(c)
			}

			return utils.NewUnauthorizedError("Authentication failed")
		}
	}
}
