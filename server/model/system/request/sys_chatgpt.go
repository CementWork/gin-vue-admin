package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type ChatGptRequest struct {
	system.Datasource
	system.ChatGpt
	request.PageInfo
}
