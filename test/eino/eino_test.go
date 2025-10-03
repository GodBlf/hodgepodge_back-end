package eino

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sync"
	"testing"
	"xmu_roll_call/app_plus/eino_agent/chatmodel"
	"xmu_roll_call/global"
	"xmu_roll_call/initialize"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	rmilvus "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"
)

func TestHello(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	ctx := context.Background()
	chatSession, err := chatmodel.NewChatSessionImpl()
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

// 测试数据
var docs = []*schema.Document{
	{
		ID:      "1",
		Content: "我是原神高手",
		MetaData: map[string]any{
			"author": "张三",
		},
	},
	{
		ID:      "2",
		Content: "你好啊",
		MetaData: map[string]any{
			"author": "李四",
		},
	},
	{
		ID:      "3",
		Content: "今天的天气真好，适合出去散步",
		MetaData: map[string]any{
			"author": "王五",
		},
	},
	{
		ID:      "4",
		Content: "我在学习Go语言编程",
		MetaData: map[string]any{
			"author": "赵六",
		},
	},
	{
		ID:      "5",
		Content: "昨晚看了一部很感人的电影",
		MetaData: map[string]any{
			"author": "钱七",
		},
	},
	{
		ID:      "6",
		Content: "我喜欢打篮球",
		MetaData: map[string]any{
			"author": "孙八",
		},
	},
	{
		ID:      "7",
		Content: "周末准备去爬山",
		MetaData: map[string]any{
			"author": "周九",
		},
	},
	{
		ID:      "8",
		Content: "最近在研究人工智能",
		MetaData: map[string]any{
			"author": "吴十",
		},
	},
	{
		ID:      "9",
		Content: "我养了一只可爱的猫",
		MetaData: map[string]any{
			"author": "郑十一",
		},
	},
	{
		ID:      "10",
		Content: "这本书非常有趣",
		MetaData: map[string]any{
			"author": "王十二",
		},
	},
	{
		ID:      "11",
		Content: "正在做一个新项目",
		MetaData: map[string]any{
			"author": "陈十三",
		},
	},
	{
		ID:      "12",
		Content: "我在学做菜",
		MetaData: map[string]any{
			"author": "刘十四",
		},
	},
	{
		ID:      "13",
		Content: "早起锻炼身体很有好处",
		MetaData: map[string]any{
			"author": "黄十五",
		},
	},
	{
		ID:      "14",
		Content: "我打算去旅游",
		MetaData: map[string]any{
			"author": "宋十六",
		},
	},
	{
		ID:      "15",
		Content: "刚买了一台新电脑",
		MetaData: map[string]any{
			"author": "方十七",
		},
	},
	{
		ID:      "16",
		Content: "我喜欢听音乐",
		MetaData: map[string]any{
			"author": "冯十八",
		},
	},
	{
		ID:      "17",
		Content: "这道菜味道真不错",
		MetaData: map[string]any{
			"author": "邓十九",
		},
	},
	{
		ID:      "18",
		Content: "我在学画画",
		MetaData: map[string]any{
			"author": "何二十",
		},
	},
	{
		ID:      "19",
		Content: "学会了新的编程技巧",
		MetaData: map[string]any{
			"author": "吕二十一",
		},
	},
	{
		ID:      "20",
		Content: "我喜欢喝咖啡",
		MetaData: map[string]any{
			"author": "施二十二",
		},
	},
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
	collection := "test"
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
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		zap.L().Error("failed to store documents", zap.Error(err))
		return
	}
	zap.L().Info("Stored document IDs", zap.Any("ids", ids))
}

func TestRetrieve(t *testing.T) {
	initialize.InitLogger()
	initialize.InitConfig()
	ctx := context.Background()
	embedder, err := openai.NewEmbedder(ctx,
		&openai.EmbeddingConfig{
			APIKey:  global.EmbedModelVar.ApiKey,
			Model:   global.EmbedModelVar.ModelName,
			BaseURL: global.EmbedModelVar.Url,
		},
	)
	if err != nil {
		zap.L().Error("failed to create embedder", zap.Error(err))
		return
	}
	milvusClient, err := client.NewClient(ctx, client.Config{
		Address:  "172.18.131.29:19530",
		Username: "godblf",
		Password: "asd456",
		DBName:   "awesomeEino",
	})
	if err != nil {
		zap.L().Error("failed to create milvus client", zap.Error(err))
		return
	}
	retriever, err := rmilvus.NewRetriever(ctx, &rmilvus.RetrieverConfig{
		Client:      milvusClient,
		Collection:  "test",
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		//返回的最大结果数目
		TopK:      8,
		Embedding: embedder,
	})
	if err != nil {
		zap.L().Error("failed to create retriever", zap.Error(err))
		return
	}
	result, err := retriever.Retrieve(ctx, "我在学习")
	if err != nil {
		zap.L().Error("failed to retrieve documents", zap.Error(err))
		return
	}
	for _, ele := range result {
		zap.L().Info("result element",
			zap.String("id", ele.ID),
			zap.String("content", ele.Content),
			zap.String("author", ele.MetaData["author"].(string)),
		)
	}
}
