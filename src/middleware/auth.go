package middleware

import (
	"net/http"

	"restApi-GoGin/src/repository"
	"restApi-GoGin/src/utils"

	"github.com/gin-gonic/gin"
)

func Auth(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("accessToken")
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token in cookie"})
			c.Abort()
			return
		}

		claims, err := utils.VerifyAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		user, err := authRepo.GetUserById(claims.UserId)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		if user.DeletedAt != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User account has been deleted"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func AuthAccess(authRepo repository.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("accessToken")
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid token in cookie"})
			c.Abort()
			return
		}

		claims, err := utils.VerifyAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		user, err := authRepo.GetUserById(claims.UserId)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		if user.DeletedAt != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User account has been deleted"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "Access denied: admin only"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
