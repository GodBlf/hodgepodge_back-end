package model

type ConfigModel struct {
	UserName          string
	PassWord          string
	UserAgent         string
	IdsUrl            string
	RollCallStatusUrl string
}
type RollCalls struct {
	List []RollCallItem `json:"rollcalls"`
}
type RollCallItem struct {
	CourseTitle string `json:"course_title"`
	CourseId    int    `json:"course_id"`
	RollcallID  int    `json:"rollcall_id"`
	IsNumber    bool   `json:"is_number"`
	IsRadar     bool   `json:"is_radar"`
	Status      string `json:"status"`
}
type RollCallJson struct {
	RollcallStatus string `json:"rollcall_status"`
	Status         string `json:"status"`
	IsExpired      bool   `json:"is_expired"`
	CourseTitle    string `json:"course_title"`
	RollcallID     int    `json:"rollcall_id"`
}

type Location struct {
	Name      string
	Longitude float64
	Latitude  float64
}

//type RollCalls struct {
//	List []RollCallJson `json:"rollcalls"`
//}

type RandomSentenceUrl struct {
	Urls []string
}

// url and api_key
type LlmUrlKey struct {
	Url    string `mapstructure:"url"`
	ApiKey string `mapstructure:"api_key"`
}

// 仅包含主力模型
type ModelLlmUrlKey struct {
	ModelCount    int         `mapstructure:"model_count"`
	Gpt5Chat      []LlmUrlKey `mapstructure:"gpt_5_chat"`
	Claude4Sonnet []LlmUrlKey `mapstructure:"claude_4_sonnet"`
}
