package eino

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sync"
	"testing"
	"xmu_roll_call/app_plus/einochat"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"
)

func TestHello(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	ctx := context.Background()
	chatSession, err := einochat.NewChatSessionImpl()
	if err != nil {
		zap.L().Error("failed to create chat session", zap.Error(err))
		return
	}
	chatSession.SetPromptTemplate("你是一个{{.role}},你的任务是{{.task}}", "{{.question}}")
	chatSession.SetSystemPrompt(ctx, map[string]any{
		"role": "大学golang老师",
		"task": "帮助学生解决学习中的问题",
	})
	reply, err := chatSession.Chat(ctx,
		map[string]any{
			"question": "讲解golang的函数闭包	",
		},
	)
	if err != nil {
		zap.L().Error("Chat error", zap.Error(err))
		return
	}
	zap.L().Info("Chat reply", zap.String("reply", reply))
	//stream midwear
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatSession.ChatStream(ctx, map[string]any{
			"question": "讲解golang的接口	",
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		chatSession.ChatStreamOut()
	}()
	wg.Wait()
}

func TestTemplate(t *testing.T) {
	tmpl := "你是{{.name}}"
	parse, err := template.New("my_template").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	//buffer字节缓冲区相当于strings.builder容器
	bf := &bytes.Buffer{}
	parse.Execute(bf, map[string]any{
		"name": "大学老师",
	})
	fmt.Println(bf.String())
	m := map[string]any{}
	v, ok := m["name"]
	if !ok {
		fmt.Println("key not found")
	} else {
		fmt.Println(v)
	}

}

func TestRag(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	ctx := context.Background()
	fmt.Println(ctx)
	zap.L().Info("embed config", zap.Any("config", global.EmbedModelVar))
	texts := []string{
		"你好",
		"我很好",
		"你是谁",
	}
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		APIKey:  global.EmbedModelVar.ApiKey,
		Model:   global.EmbedModelVar.ModelName,
		BaseURL: global.EmbedModelVar.Url,
	})
	if err != nil {
		zap.L().Error("failed to create embedder", zap.Error(err))
		return
	}
	result, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		zap.L().Error("failed to embed strings", zap.Error(err))
		return
	}
	zap.L().Info("Embed result", zap.Any("result dimension", len(result[0])))
	//indexer
	milclient, err := client.NewClient(ctx, client.Config{
		Address:  "172.18.131.29:19530",
		Username: "godblf",
		Password: "asd456",
		DBName:   "awesomeEino",
	})
	if err != nil {
		zap.L().Error("failed to create milvus client", zap.Error(err))
		return
	}
	collection := "test1"
	fields := []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "1024",
			},
			PrimaryKey: true,
		},
		{
			Name: "vector",
			//todo:垃圾阿到底32还是64
			//向量数据库行业标准32位,64内部自动转成32位阿!!!!
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "32768",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "1024",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     milclient,
		Collection: collection,
		Fields:     fields,
		Embedding:  embedder,
	})
	if err != nil {
		zap.L().Error("failed to create milvus indexer", zap.Error(err))
		return
	}
	docs := []*schema.Document{
		{
			ID:      "1",
			Content: "你好",
			MetaData: map[string]any{
				"author": "鲁迅",
			},
		},
		{
			ID:      "2",
			Content: "我很好",
			MetaData: map[string]any{
				"author": "诸葛亮",
			},
		},
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		zap.L().Error("failed to store documents", zap.Error(err))
		return
	}
	zap.L().Info("Stored document IDs", zap.Any("ids", ids))
}
