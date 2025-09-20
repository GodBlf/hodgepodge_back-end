package roll_call

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"xmu_roll_call/app/login"
	"xmu_roll_call/global"
	"xmu_roll_call/model"
	"xmu_roll_call/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

var (
	RollCallImplVar *RollCallImpl
	once            sync.Once
)

type RollCallImpl struct {
	Login    login.Login
	Client   *resty.Client
	DeviceId string
}

//	func (r *RollCallImpl) RollCallFinal(c *gin.Context) {
//		courseAndRollcallId, err := r.RollCallStatus()
//		if err != nil {
//			zap.L().Error("查询签到状态失败", zap.Error(err))
//			c.JSON(500, gin.H{"message": "查询签到状态失败"})
//			return
//		}
//		zap.L().Info("查询签到状态成功", zap.Any("查询结果", courseAndRollcallId))
//		if len(courseAndRollcallId) == 0 {
//			zap.L().Info("当前没有需要签到的课程")
//			c.JSON(200, gin.H{"message": "当前没有需要签到的课程"})
//			return
//		}
//		zap.L().Info("开始查询签到码...")
//		query, err, rollcallType := r.NumberCodeQuery(courseAndRollcallId)
//
//		if rollcallType == 0 {
//			zap.L().Info("签到码查询成功", zap.Any("查询结果", query))
//			zap.L().Info("打印签到码...")
//			for title, numbercode := range query {
//				if numbercode != "" {
//					fmt.Printf("✅ %s: 签到码 %s\n", title, numbercode)
//				} else {
//					fmt.Printf("❌ %s: 获取签到码失败\n", title)
//				}
//			}
//			zap.L().Info("开启自动签到...")
//
//			err = r.NumberCodePost(courseAndRollcallId, query, r.DeviceId)
//			if err != nil {
//				zap.L().Error("自动签到失败", zap.Error(err))
//				query["message"] = "自动签到失败"
//				c.JSON(500, query)
//				return
//			}
//			zap.L().Info("自动签到成功")
//			query["message"] = "自动签到成功"
//			c.JSON(200, query)
//			return
//
//		} else if rollcallType == 1 {
//			//radar rollcall
//			zap.L().Info("雷达签到", zap.Any("查询结果", query))
//			zap.L().Info("开启自动签到...")
//			c.JSON(200, gin.H{"message": "雷达签到"})
//			return
//		} else {
//			zap.L().Warn("二维码签到请老实上课", zap.Any("查询结果", query))
//			c.JSON(200, gin.H{"message": "二维码签到请老实上课"})
//			return
//		}
//	}
func (r *RollCallImpl) RollCallFinal(c *gin.Context) {
	courses, err := r.RollCallStatus()
	if err != nil {
		zap.L().Error("查询签到状态失败", zap.Error(err))
		c.JSON(500, gin.H{"message": "查询签到状态失败"})
		return
	}
	if len(courses.List) == 0 {
		c.JSON(200, gin.H{"message": "当前没有需要签到的课程"})
		return
	}

	results := make([]string, 0)

	for _, course := range courses.List {
		if course.Status != "absent" {
			results = append(results, fmt.Sprintf("✅ %s 已完成签到", course.CourseTitle))
			continue
		}

		// 数字签到
		if course.IsNumber {
			numberCode, err := r.getNumberCode(course.RollcallID)
			if err != nil || numberCode == "" {
				results = append(results, fmt.Sprintf("❌ %s 获取签到码失败", course.CourseTitle))
				continue
			}
			err = r.postNumberCode(course.RollcallID, numberCode, r.DeviceId)
			if err != nil {
				results = append(results, fmt.Sprintf("❌ %s 数字签到失败", course.CourseTitle))
				continue
			}
			results = append(results, fmt.Sprintf("✅ %s 数字签到成功，签到码 %s", course.CourseTitle, numberCode))
			continue
		}

		// 雷达签到
		if course.IsRadar {
			lat, lon, locName, dist, err := r.getRadarBestLocation(course.RollcallID)
			if err != nil {
				results = append(results, fmt.Sprintf("❌ %s 雷达定位失败", course.CourseTitle))
				continue
			}
			err = r.postRadarSign(course.RollcallID, lat, lon, r.DeviceId)
			if err != nil {
				results = append(results, fmt.Sprintf("❌ %s 雷达签到失败", course.CourseTitle))
				continue
			}
			results = append(results, fmt.Sprintf("✅ %s 雷达签到成功，位置 %s 距离 %.2fm", course.CourseTitle, locName, dist))
			continue
		}

		// 二维码签到
		results = append(results, fmt.Sprintf("⚠️ %s 二维码签到，请人工完成", course.CourseTitle))
	}

	c.JSON(200, gin.H{"result": results})
}
func (r *RollCallImpl) getNumberCode(rollcallID int) (string, error) {
	url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/student_rollcalls", rollcallID)
	resp, err := r.Client.R().Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() >= 400 {
		return "", fmt.Errorf("status %d", resp.StatusCode())
	}
	code := gjson.Get(resp.String(), "number_code").String()
	return code, nil
}

func (r *RollCallImpl) postNumberCode(rollcallID int, numberCode, deviceId string) error {
	url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/answer_number_rollcall", rollcallID)
	payload := map[string]string{
		"deviceId":   deviceId,
		"numberCode": numberCode,
	}
	resp, err := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
}

type RadarSign struct {
	Distance float64 `json:"distance"`
}

func (r *RollCallImpl) getRadarBestLocation(rollcallID int) (lat, lon float64, name string, dist float64, err error) {
	bestDist := 1e9
	var bestLoc model.Location

	for _, loc := range global.Locations {
		payload := map[string]interface{}{
			"deviceId":         r.DeviceId,
			"latitude":         loc.Latitude,
			"longitude":        loc.Longitude,
			"speed":            nil,
			"accuracy":         90,
			"altitude":         nil,
			"altitudeAccuracy": nil,
			"heading":          nil,
		}
		url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/answer?api_version=1.1.2", rollcallID)
		resp, err := r.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(payload).
			Put(url)
		if err != nil || resp.StatusCode() >= 400 {
			continue
		}

		var rs RadarSign
		if e := json.Unmarshal(resp.Body(), &rs); e != nil {
			continue
		}

		if rs.Distance < bestDist {
			bestDist = rs.Distance
			bestLoc = loc
			dist = rs.Distance
		}
		if rs.Distance < 100 {
			break
		}
	}

	if bestDist == 1e9 {
		return 0, 0, "", 0, fmt.Errorf("无法找到合适位置")
	}

	// 抖动
	lat = bestLoc.Latitude + utils.RandomRange(-0.0003, 0.0003)
	lon = bestLoc.Longitude + utils.RandomRange(-0.0003, 0.0003)
	name = bestLoc.Name
	return lat, lon, name, dist, nil
}

func (r *RollCallImpl) postRadarSign(rollcallID int, lat, lon float64, deviceId string) error {
	url := fmt.Sprintf("https://lnt.xmu.edu.cn/api/rollcall/%d/answer?api_version=1.1.2", rollcallID)
	payload := map[string]interface{}{
		"deviceId":         deviceId,
		"latitude":         lat,
		"longitude":        lon,
		"speed":            nil,
		"accuracy":         90,
		"altitude":         nil,
		"altitudeAccuracy": nil,
		"heading":          nil,
	}
	resp, err := r.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Put(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
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
	zap.L().Info("登录成功", zap.String("deviceID", deviceID))
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
			Client:   c,
			Login:    l,
			DeviceId: utils.Uuid(),
		}
	})
	return RollCallImplVar
}
func (r *RollCallImpl) RollCallStatus() (*model.RollCalls, error) {
	url := global.Config.RollCallStatusUrl
	response, err := r.Client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() >= 400 {
		return nil, fmt.Errorf("status %d", response.StatusCode())
	}
	var unmarshalList model.RollCalls
	if err := json.Unmarshal(response.Body(), &unmarshalList); err != nil {
		return nil, err
	}
	return &unmarshalList, nil
}

//func (r *RollCallImpl) RollCallStatus() (map[string]int, error) {
//	url := global.Config.RollCallStatusUrl
//	response, err := r.Client.R().Get(url)
//	if err != nil {
//		zap.L().Error("签到状态请求失败", zap.Error(err))
//		return nil, err
//	}
//	if response.StatusCode() >= 400 {
//		zap.L().Error("签到状态请求返回错误状态码", zap.Int("status_code", response.StatusCode()))
//		return nil, err
//	}
//	unmarshalList := &model.RollCalls{}
//	err = json.Unmarshal(response.Body(), unmarshalList)
//	if err != nil {
//		zap.L().Error("签到状态响应解析失败", zap.Error(err))
//		return nil, err
//	}
//	pending := make(map[string]int)
//	for _, status := range unmarshalList.List {
//		if status.RollcallID != 0 {
//			pending[status.CourseTitle] = status.RollcallID
//		}
//	}
//	return pending, nil
//}
