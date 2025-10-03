package eino_milvus

import (
	"context"
	"sync"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var (
	MilvusClientVar client.Client
	monce           = &sync.Once{}
)

func NewMilvusClient(ctx context.Context) client.Client {
	monce.Do(func() {
		milclient, err := client.NewClient(ctx, client.Config{
			Address:  "172.18.131.29:19530",
			Username: "godblf",
			Password: "asd456",
			DBName:   "awesomeEino",
		})
		if err != nil {
			panic(err)
		}
		MilvusClientVar = milclient
	})
	return MilvusClientVar
}
