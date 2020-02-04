package handler

import "github.com/gin-gonic/gin"

func handleError(args ...error) gin.H {
	var errors []gin.H
	for _, arg := range args {
		errors = append(errors, gin.H{"message": arg.Error()})
	}
	return gin.H{"errors": errors}
}
