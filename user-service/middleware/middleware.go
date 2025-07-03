package middleware

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jstnangrendo/instagram-clone/user-service/config"
	"github.com/jstnangrendo/instagram-clone/user-service/domains/users/entities"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing Authorization header"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		userIDf, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}
		c.Set("user_id", uint(userIDf))

		tid, _ := claims["token_id"].(string)
		var at entities.AccessToken

		if raw, err := config.RedisClient.Get(config.Ctx, "access_token:"+tid).Result(); err == nil {
			_ = json.Unmarshal([]byte(raw), &at)
		} else {
			if err := config.DB.First(&at, "id = ?", tid).Error; err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "token not found"})
				return
			}
			buf, _ := json.Marshal(at)
			config.RedisClient.Set(config.Ctx, "access_token:"+tid, buf, time.Until(at.ExpiresAt))
		}

		if at.Revoked || time.Now().After(at.ExpiresAt) {
			c.AbortWithStatusJSON(401, gin.H{"error": "token revoked or expired"})
			return
		}
		c.Next()
	}
}
