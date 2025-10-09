package eino_embedder

import (
	"context"
	"sync"
	"xmu_roll_call/global"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

var (
	EmbedderVar *openai.Embedder
	eonce       = &sync.Once{}
)

func NewEmbedderBge() *openai.Embedder {
	eonce.Do(func() {
		e, err := openai.NewEmbedder(context.Background(), &openai.EmbeddingConfig{
			APIKey:  global.EmbedModelVar.ApiKey,
			Model:   global.EmbedModelVar.ModelName,
			BaseURL: global.EmbedModelVar.Url,
		})
		if err != nil {
			panic(err)
		}
		EmbedderVar = e
	})

	return EmbedderVar
}
