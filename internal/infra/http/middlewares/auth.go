package middlewares

import (
	"net/http"
	"strings"

	"apsdigital/internal/config"
	"apsdigital/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Level  int       `json:"level"`
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_level", claims.Level)

		c.Next()
	}
}

func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		// Check if user has required role
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func RequireMinLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, exists := c.Get("user_level")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User level not found"})
			c.Abort()
			return
		}

		level, ok := userLevel.(int)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user level"})
			c.Abort()
			return
		}

		// Lower level number means higher hierarchy (1=ADM, 2=Coordenador, etc.)
		if level > minLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CanAuthorize checks if current user can authorize another user based on hierarchy
func CanAuthorize(targetUserLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel, exists := c.Get("user_level")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User level not found"})
			c.Abort()
			return
		}

		level, ok := userLevel.(int)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user level"})
			c.Abort()
			return
		}

		// Check authorization hierarchy:
		// ADM (1) can authorize anyone
		// Coordenador (2) can authorize Gerente (3) and ACS (4)
		// Gerente (3) can authorize ACS (4)
		// ACS (4) cannot authorize anyone

		canAuth := false
		switch level {
		case entities.LevelAdmin: // ADM
			canAuth = true
		case entities.LevelCoordenador: // Coordenador
			canAuth = targetUserLevel >= entities.LevelGerente
		case entities.LevelGerente: // Gerente
			canAuth = targetUserLevel == entities.LevelACS
		case entities.LevelACS: // ACS
			canAuth = false
		}

		if !canAuth {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot authorize users of this level"})
			c.Abort()
			return
		}

		c.Next()
	}
}