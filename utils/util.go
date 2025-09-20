package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
	"xmu_roll_call/global"
	"xmu_roll_call/model"
)

// todo:可以直接在 utils 包里放所有工具函数： 等项目成长到一定规模再重构。
func Exe2DeviceID(execution string) (string, error) {
	index := strings.Index(execution, "_")
	if index == -1 {
		return "", errors.New("execution格式错误")
	}
	return execution[0:index], nil
}

// todo:最短路径
func MinDistance() (string, error) {
	return "", nil
}

// 随机返回一个模型的url和apikey
func RandomLlmModel() (*model.LlmUrlKey, string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	choose := r.Intn(global.MLUK.ModelCount)
	if choose == 0 {
		if len(global.MLUK.Gpt5Chat) == 0 {
			return nil, "", errors.New("没有可用的模型")
		}
		tmp := r.Intn(len(global.MLUK.Gpt5Chat))
		return &global.MLUK.Gpt5Chat[tmp], "gpt-5-chat", nil
	} else {
		if len(global.MLUK.Claude4Sonnet) == 0 {
			return nil, "", errors.New("没有可用的模型")
		}
		tmp := r.Intn(len(global.MLUK.Claude4Sonnet))
		return &global.MLUK.Claude4Sonnet[tmp], "claude-4-sonnet", nil
	}

}

// 随机返回一个指定的模型url和apikey
func RandomSpecifiedLlmModel(modelName string) (*model.LlmUrlKey, error) {
	var models []model.LlmUrlKey
	switch modelName {
	case "gpt_5_chat":
		models = global.MLUK.Gpt5Chat
	case "claude_4_sonnet":
		models = global.MLUK.Claude4Sonnet
	default:
		return nil, errors.New("不支持的模型名称")
	}
	if len(models) == 0 {
		return nil, errors.New("没有可用的模型")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(models))
	return &models[idx], nil
}
