package eino_indexer

import (
	"context"
	"sync"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var (
	MilvusIndexerVar *MilvusIndexerImpl
	milonce          = &sync.Once{}
)

type MilvusIndexerImpl struct {
	Embedder     *openai.Embedder
	MilvusClient client.Client
	Indexer      *milvus.Indexer
}

func NewMilvusClient(ctx context.Context, e *openai.Embedder, m client.Client) *MilvusIndexerImpl {
	milonce.Do(func() {
		milvus.NewIndexer(ctx, &milvus.IndexerConfig{})
		mii := &MilvusIndexerImpl{
			Embedder:     e,
			MilvusClient: m,
		}
		MilvusIndexerVar = mii
	})
	return MilvusIndexerVar
}
