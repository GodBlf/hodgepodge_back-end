package roll_call

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"sync"
	"xmu_roll_call/app/login"
	"xmu_roll_call/global"
	"xmu_roll_call/model"
	"xmu_roll_call/utils"
)

var (
	RollCallImplVar *RollCallImpl
	once            sync.Once
)

type RollCallImpl struct {
	Login  login.Login
	Client *resty.Client
}

func (r *RollCallImpl) RollCallLogin() (string, error) {
	execution, boolean, err := r.Login.Login(global.Config.UserName, global.Config.PassWord)
	if err != nil {
		zap.L().Error("登录失败", zap.Error(err))
		return execution, err
	}
	deviceID, err := utils.Exe2DeviceID(execution)
	if err != nil {
		zap.L().Error("设备ID生成失败", zap.Error(err))
		return "", err
	}

	if !boolean {
		zap.L().Warn("登录失败,请检查用户名和密码是否正确", zap.String("execution", execution))
		return deviceID, errors.New("登录失败,请检查用户名和密码是否正确")
	}
	zap.L().Info("登录成功", zap.String("execution", execution))
	return deviceID, nil

}

//0:number 1:radar 2:qr
//todo:以后用enum枚举代替
//var RadarError = errors.New("radar_rollcall")

func (r *RollCallImpl) NumberCodeQuery(rollcall map[string]int) (map[string]string, error, int) {
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
		responseString := response.String()
		isRadar := gjson.Get(responseString, "is_radar").Bool()
		if isRadar {
			return nil, nil, 1
		}
		NumberCode := gjson.Get(responseString, "number_code").String()
		if NumberCode != "" {
			results[title] = NumberCode
			zap.L().Info("课程签到码查询成功", zap.String("course", title), zap.String("number_code", NumberCode))
		} else {
			zap.L().Warn("课程未找到签到码", zap.String("course", title))
			results[title] = ""
		}
	}
	return results, nil, 0
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

func NewRollCallImpl(c *resty.Client, l login.Login) *RollCallImpl {
	once.Do(func() {
		RollCallImplVar = &RollCallImpl{
			Client: c,
			Login:  l,
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
