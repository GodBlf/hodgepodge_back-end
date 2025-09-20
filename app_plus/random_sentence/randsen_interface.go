package random_sentence

import "github.com/gin-gonic/gin"

type RandSenInterface interface {
	send(context *gin.Context)
	get(input string) (string, error)
}
