package global

import (
	"xmu_roll_call/model"

	"github.com/gin-gonic/gin"
)

var (
	Config *model.ConfigModel
	Router *gin.Engine
)

var EmbedModelVar *model.EmbedModel = &model.EmbedModel{}

var (
	LlmUrl        string
	LlmApiKey     string
	LlmModel      string
	OpenaiVersion string
)

var Locations = []model.Location{
	{"学武楼", 118.313801, 24.605586},
	{"文宣楼", 118.309978, 24.605288},
	{"新工科大楼", 118.310252, 24.614584},
	{"坤銮楼", 118.312747, 24.605544},
	{"西部片区2号楼", 118.299847, 24.604192},
	{"西部片区4号楼", 118.30018608783871, 24.60527088060157},
	{"一期田径场", 118.31887190174575, 24.608956831402885},
	{"游泳馆", 118.31191862035234, 24.610805475596482},
	{"爱秋体育馆", 118.31051001664491, 24.61151956997547},
	{"一期篮球场", 118.31723671918621, 24.608388221582917},
}

var MLUK *model.ModelLlmUrlKey
