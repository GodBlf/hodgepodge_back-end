package model

type ConfigModel struct {
	UserName          string
	PassWord          string
	UserAgent         string
	IdsUrl            string
	RollCallStatusUrl string
}

type RollCallJson struct {
	RollcallStatus string `json:"rollcall_status"`
	Status         string `json:"status"`
	IsExpired      bool   `json:"is_expired"`
	CourseTitle    string `json:"course_title"`
	RollcallID     int    `json:"rollcall_id"`
}
