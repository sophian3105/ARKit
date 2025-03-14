package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	HandleRouter(router)
	router.Run("localhost:8080")
}