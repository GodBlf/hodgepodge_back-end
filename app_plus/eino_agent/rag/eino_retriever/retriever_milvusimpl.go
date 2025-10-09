package eino_retriever

import (
	"context"
	"sync"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	rmilvus "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

var (
	RetrieverMilvusVar *RetrieverMilvusRam
	ronce              = &sync.Once{}
)

type RetrieverMilvusRam struct {
	MilvusClientP *client.Client
	EmbedderBgeP  *openai.Embedder
	RetrieverP    *rmilvus.Retriever
}

func NewRetrieverMilvus(m *client.Client, e *openai.Embedder) *RetrieverMilvusRam {
	ronce.Do(func() {
		rtmp, err := rmilvus.NewRetriever(context.Background(), &rmilvus.RetrieverConfig{
			TopK:       4,
			Client:     *m,
			Collection: "test",
			OutputFields: []string{
				"id",
				"content",
				"metadata",
			},
			Embedding: e,
		})
		if err != nil {
			zap.L().Error("failed to create milvus retriever", zap.Error(err))
			return
		}
		rmi := &RetrieverMilvusRam{
			MilvusClientP: m,
			EmbedderBgeP:  e,
			RetrieverP:    rtmp,
		}
		RetrieverMilvusVar = rmi
	})
	return RetrieverMilvusVar
}
