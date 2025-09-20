package utils

import (
	"fmt"
	"testing"
	"xmu_roll_call/utils"
)

func TestUuid(t *testing.T) {
	uuid := utils.Uuid()
	fmt.Println(uuid)
}
