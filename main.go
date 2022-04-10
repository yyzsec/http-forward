package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	err := r.Run("0.0.0.0:8081")
	if err != nil {
		return 
	}
}