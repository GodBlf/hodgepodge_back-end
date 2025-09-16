package llms

import (
	"testing"
	"xmu_roll_call/app_plus/llms"
	"xmu_roll_call/initialize"
)

func TestTranslate(t *testing.T) {
	initialize.InitConfig()
	resp, err := llms.Translate("你好")
	if err != nil {
		t.Errorf("Translate failed: %v", err)
		return
	}
	t.Logf("Translate response: %s", resp)
}
