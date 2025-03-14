package main

import (
	util "aria/backend/utility"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// HandleRouter provides the paths for the router
// Owner must call router.Run(...)
func HandleRouter(router *gin.Engine) {
	router.GET("/ping", getPing)

	// Any routes/groups in authorized must be secured 
	// with firebase
	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/auth", getAuth)
	}
}

// Simple get method to ping the backend
func getPing(c *gin.Context) {
	ginContent := gin.H{
		"message": "pong",
	}

	c.JSON(http.StatusOK, ginContent)
}

// Simple method to test auth
func getAuth(c *gin.Context) {
	token := c.MustGet("token").(auth.Token)
	c.JSON(http.StatusOK, token)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := util.FBClient()
		
		if client == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		authHeader := c.GetHeader("Authorization")
		authTokenString, found := strings.CutPrefix(authHeader, "Bearer ")

		if !found {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken, err := client.Auth.VerifyIDToken(c, authTokenString)
		
		if err != nil || authToken != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("token", authToken)
		c.Next()
	}
}
