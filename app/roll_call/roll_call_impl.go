package roll_call

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"sync"
	"xmu_roll_call/global"
	"xmu_roll_call/model"
)

var (
	RollCallImplVar *RollCallImpl
	once            sync.Once
)

type RollCallImpl struct {
	Client *resty.Client
}

func (r *RollCallImpl) NumberCodeQuery(rollcall map[string]int) (map[string]string, error) {
	//TODO implement me
	results := make(map[string]string)
	for title, rollCallId := range rollcall {
		url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/student_rollcalls", rollCallId)
		response, err := r.Client.R().Get(url)
		if err != nil {
			zap.L().Error("签到码查询请求失败", zap.Error(err))
			results[title] = ""
			continue
		}
		if response.StatusCode() >= 400 {
			zap.L().Error("签到码查询请求返回错误状态码", zap.Int("status_code", response.StatusCode()))
			results[title] = ""
			continue
		}

		if err != nil {
			zap.L().Error("签到码查询响应解析失败", zap.Error(err))
			results[title] = ""
			continue
		}
		responseString := response.String()
		NumberCode := gjson.Get(responseString, "number_code").String()
		if NumberCode != "" {
			results[title] = NumberCode
			zap.L().Info("课程签到码查询成功", zap.String("course", title), zap.String("number_code", NumberCode))
		} else {
			zap.L().Warn("课程未找到签到码", zap.String("course", title))
			results[title] = ""
		}
	}
	return results, nil
}

func (r *RollCallImpl) NumberCodePost(courseNameRollCallId map[string]int, numberCode map[string]string, deviceId string) error {
	for courseName, code := range numberCode {
		url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/answer_number_rollcall", courseNameRollCallId[courseName])
		payload := map[string]string{
			"deviceId":   deviceId,
			"numberCode": code,
		}
		response, err := r.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Post(url)
		if err != nil {
			zap.L().Error("签到请求失败", zap.String("course", courseName), zap.Error(err))
			return err
		}
		if response.StatusCode() >= 400 {
			zap.L().Error("签到请求返回错误状态码", zap.String("course", courseName), zap.Int("status_code", response.StatusCode()), zap.String("response", response.String()))
			return err
		}
		zap.L().Info("课程签到成功", zap.String("course", courseName), zap.String("response", response.String()))
	}
	return nil
}

func NewRollCallImpl(c *resty.Client) *RollCallImpl {
	once.Do(func() {
		RollCallImplVar = &RollCallImpl{
			Client: c,
		}
	})
	return RollCallImplVar
}

func (r *RollCallImpl) RollCallStatus() (map[string]int, error) {
	url := global.Config.RollCallStatusUrl
	response, err := r.Client.R().Get(url)
	if err != nil {
		zap.L().Error("签到状态请求失败", zap.Error(err))
		return nil, err
	}
	if response.StatusCode() >= 400 {
		zap.L().Error("签到状态请求返回错误状态码", zap.Int("status_code", response.StatusCode()))
		return nil, err
	}
	unmarshalList := struct{ list []model.RollCallJson }{list: make([]model.RollCallJson, 10)}
	err = json.Unmarshal(response.Body(), &unmarshalList)
	if err != nil {
		zap.L().Error("签到状态响应解析失败", zap.Error(err))
		return nil, err
	}
	pending := make(map[string]int)
	for _, status := range unmarshalList.list {
		if status.RollcallID != 0 {
			pending[status.CourseTitle] = status.RollcallID
		}
	}
	return pending, nil
}
