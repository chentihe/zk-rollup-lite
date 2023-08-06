package controllers

import "github.com/gin-gonic/gin"

// TODO: read how to use gin.NoMethod()
func HandleError(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{"error": err.Error()})
}
