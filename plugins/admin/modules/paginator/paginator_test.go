package paginator

import (
	"testing"

	"goadmin/modules/config"
	"goadmin/plugins/admin/modules/parameter"
	_ "goadmin/themes/sword"
)

func TestGet(t *testing.T) {
	config.Initialize(&config.Config{Theme: "sword"})
	param := parameter.BaseParam()
	param.Page = "7"
	Get(nil, Config{
		Size:         105,
		Param:        param,
		PageSizeList: []string{"10", "20", "50", "100"},
	})
}
