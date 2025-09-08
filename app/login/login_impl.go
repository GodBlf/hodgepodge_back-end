package login

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"strings"
	"sync"
	"xmu_roll_call/app/encrypt"
	"xmu_roll_call/global"
)

var (
	LoginVar *LoginImpl
	once     sync.Once
)

type LoginImpl struct {
	Encrypt encrypt.Encrypt
	Client  *resty.Client
}

func (l *LoginImpl) GetLoginPage() (salt, execution, lt string, err error) {
	url := global.Config.IdsUrl
	resp, err := l.Client.R().Get(url)
	if err != nil {
		zap.L().Error("请求登录页失败", zap.Error(err))
		return "", "", "", err
	}
	if resp.StatusCode() >= 400 {
		zap.L().Error("请求登录页返回错误状态码", zap.Int("status_code", resp.StatusCode()))
		return "", "", "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		zap.L().Error("解析登录页失败", zap.Error(err))
		return "", "", "", err
	}
	salt, boolean := doc.Find("#pwdEncryptSalt").Attr("value")
	if boolean == false {
		zap.L().Error("加密盐抓取失败", zap.Error(err))
	} else {
		zap.L().Info("抓取加密盐成功", zap.String("encryp_salt", salt))
	}
	execution, boolean = doc.Find("input[name='execution']").Attr("value")
	if boolean == false {
		zap.L().Error("execution抓取失败", zap.Error(err))
	} else {
		zap.L().Info("execution抓取成功", zap.String("execution", execution))
	}
	lt, _ = doc.Find("input[name='lt']").Attr("value")
	if salt == "" || execution == "" {
		zap.L().Error("未从登录页提取到必要字段 salt/execution")
		return "", "", "", err
	}
	return salt, execution, lt, nil

}

func (l *LoginImpl) Login(username, password string) (string, bool, error) {
	salt, execution, lt, err := l.GetLoginPage()
	if err != nil {
		zap.L().Error("获取登录页失败", zap.Error(err))
		return execution, false, err
	}
	encPwd := l.Encrypt.EncryptPassword(password, salt)
	data := map[string]string{
		"username":  username,
		"password":  encPwd,
		"lt":        lt,
		"execution": execution,
		"_eventId":  "submit",
		"rmShown":   "1",
	}
	request := l.Client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Referer", global.Config.IdsUrl).
		SetFormData(data)
	response, err := request.Post(global.Config.IdsUrl)
	if err != nil {
		zap.L().Error("登录请求失败", zap.Error(err))
		return execution, false, err
	}
	if response.StatusCode() == 200 {
		zap.L().Info("重新定向登录成功")
		return execution, true, nil
	}
	body2string := response.String()
	if strings.Contains(body2string, "用户名或密码错误") {
		zap.L().Warn("用户名或密码错误")
		return execution, false, nil
	}
	zap.L().Warn("登录状态未知",
		zap.Int("响应码", response.StatusCode()),
		zap.String("响应体", body2string),
	)
	return execution, false, nil

}

func NewLoginImpl(client *resty.Client, encrypt encrypt.Encrypt) *LoginImpl {
	once.Do(func() {
		LoginVar = &LoginImpl{
			Encrypt: encrypt,
			Client:  client,
		}
	})
	return LoginVar
}
