package llms

import "github.com/gin-gonic/gin"

type llmInterface interface {
	Send(context *gin.Context)
}
