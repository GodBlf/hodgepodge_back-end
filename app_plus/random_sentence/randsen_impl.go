package random_sentence

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"sync"
)

var (
	RandSenImplVar *RandSenImpl
	once           = &sync.Once{}
)

type RandSenImpl struct {
	Client *resty.Client
}

func NewRandSenImpl(c *resty.Client) *RandSenImpl {
	once.Do(func() {
		RandSenImplVar = &RandSenImpl{
			Client: c,
		}
	})
	return RandSenImplVar
}

func (r *RandSenImpl) send(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r *RandSenImpl) get(input string) (string, error) {
	//TODO implement me
	panic("implement me")
}
