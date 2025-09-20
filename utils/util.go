package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"
	"xmu_roll_call/global"
	"xmu_roll_call/model"
)

func RandomRange(min, max float64) float64 {
	if min > max {
		min, max = max, min
	}
	return min + rand.Float64()*(max-min)
}
func init() {
	rand.Seed(time.Now().UnixNano())
}

// uuid
func Uuid() string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子（生产环境建议使用 crypto/rand）

	template := "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	result := make([]byte, len(template))

	for i := 0; i < len(template); i++ {
		char := template[i]
		if char != 'x' && char != 'y' && char != '-' && char != '4' {
			continue // 不处理其他字符（理论上不会出现）
		}

		if char == '-' || char == '4' {
			result[i] = char
			continue
		}

		// 生成 0-15 的随机整数
		e := rand.Intn(16)

		var r int
		if char == 'x' {
			r = e
		} else if char == 'y' {
			// y: (e & 3) | 8 → 确保高位是 10xx → 结果是 8,9,a,b
			r = (e & 3) | 8
		}

		// 转换为十六进制字符 (0-9, a-f)
		if r < 10 {
			result[i] = byte('0' + r)
		} else {
			result[i] = byte('a' + r - 10)
		}
	}

	return string(result)
}

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
