package eino_embedder

import (
	"context"
	"sync"
	"xmu_roll_call/global"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

var (
	Embeddervar *openai.Embedder
	eonce       = &sync.Once{}
)

func NewEmbedderBge(ctx context.Context) *openai.Embedder {
	eonce.Do(func() {
		e, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
			APIKey:  global.EmbedModelVar.ApiKey,
			Model:   global.EmbedModelVar.ModelName,
			BaseURL: global.EmbedModelVar.Url,
		})
		if err != nil {
			panic(err)
		}
		Embeddervar = e
	})

	return Embeddervar
}
