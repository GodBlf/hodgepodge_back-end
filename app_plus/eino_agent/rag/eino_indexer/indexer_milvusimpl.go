package eino_indexer

import (
	"context"
	"sync"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var (
	MilvusIndexerVar *MilvusIndexerRam
	milonce          = &sync.Once{}
)

type MilvusIndexerRam struct {
	EmbedderP     *openai.Embedder
	MilvusClientP *client.Client
	IndexerP      *milvus.Indexer
}

func NewMilvusClient(e *openai.Embedder, m *client.Client) *MilvusIndexerRam {
	milonce.Do(func() {
		milvus.NewIndexer(context.Background(), &milvus.IndexerConfig{})
		mii := &MilvusIndexerRam{
			EmbedderP:     e,
			MilvusClientP: m,
		}
		MilvusIndexerVar = mii
	})
	return MilvusIndexerVar
}
