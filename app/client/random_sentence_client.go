package client

import (
	"github.com/go-resty/resty/v2"
	"sync"
)

var (
	RandomSentenceClient *resty.Client
	onceRandomSentence   = &sync.Once{}
)

func NewRandomSentenceClient() *resty.Client {
	onceRandomSentence.Do(func() {
		RandomSentenceClient = resty.New()
	})
	return RandomSentenceClient
}
