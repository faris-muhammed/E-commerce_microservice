package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims struct for JWT claims
type Claims struct {
	Username string `json:"username"`
	UserId   uint
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("SECRETKEY"))

// Middleware for checking JWT validity
// func AuthMiddleware(requiredRole string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Retrieve token from Authorization header (Bearer token)
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			c.JSON(401, gin.H{
// 				"status":  "Unauthorized",
// 				"message": "No token provided",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Parse the token
// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})
// 		if err != nil || !token.Valid {
// 			c.JSON(401, gin.H{
// 				"status":  "Unauthorized",
// 				"message": "Invalid or expired JWT token",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Check for role-based authorization
// 		if claims.Role != requiredRole {
// 			c.JSON(403, gin.H{
// 				"status":  "Forbidden",
// 				"message": "Insufficient permissions",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Set user data in context for downstream handlers
// 		c.Set("userId", claims.UserId)
// 		c.Set("username", claims.Username)

//			c.Next()
//		}
//	}
// func AuthMiddleware(requiredRole string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Try to retrieve the token from the Authorization header
// 		authHeader := c.GetHeader("Authorization")
// 		var tokenString string

// 		// If no Authorization header, check for the JWT in cookies
// 		if authHeader != "" {
// 			const prefix = "Bearer "
// 			if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
// 				tokenString = authHeader[len(prefix):]
// 			}
// 		} else {
// 			// Check for the JWT in the cookie
// 			cookie, err := c.Cookie("jwtTokenAdmin")
// 			if err != nil {
// 				c.JSON(http.StatusUnauthorized, gin.H{
// 					"status":  "Unauthorized",
// 					"message": "No token provided in header or cookie",
// 				})
// 				c.Abort()
// 				return
// 			}
// 			tokenString = cookie
// 		}

// 		// Parse the token
// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})
// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"status":  "Unauthorized",
// 				"message": "Invalid or expired JWT token",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Check for role-based authorization
// 		if claims.Role != requiredRole {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"status":  "Forbidden",
// 				"message": "Insufficient permissions",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Set user data in context for downstream handlers
// 		c.Set("userId", claims.UserId)
// 		c.Set("username", claims.Username)

// 		c.Next()
// 	}
// }

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("jwtkey", jwtKey)
		tokenString, err := c.Cookie("jwtToken" + requiredRole)
		fmt.Println("TokenString", tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Can't find cookie",
				"code":    401,
			})
			c.Abort()
			return
		}
		if tokenString == "" {
			c.JSON(400, gin.H{
				"status":  "Bad Request",
				"message": "Empty token string.",
				"code":    400,
			})
			c.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Tokenclaims", token.Claims)
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			fmt.Println("cookie error:", err)
			c.JSON(401, gin.H{
				"status":  "Unauthorized",
				"message": "Invalid or expired JWT Token.",
				"code":    401,
			})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			fmt.Println("req", requiredRole, claims.Role)
			c.JSON(403, gin.H{
				"status": "Forbidden",
				"error":  "Insufficient permissions",
				"code":   403,
			})
			c.Abort()
			return
		}
		c.Set("userid", claims.UserId)
		c.Next()
	}
}
