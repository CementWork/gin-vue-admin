package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatGptApi struct{}

func (chat *ChatGptApi) CreateSK(c *gin.Context) {
	var option sysModel.SysChatGptOption
	c.ShouldBindJSON(&option)
	err := chatGptService.CreateSK(option)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (chat *ChatGptApi) GetSK(c *gin.Context) {
	var option sysModel.SysChatGptOption
	c.ShouldBindJSON(&option)
	_, err := chatGptService.GetSK()
	if err != nil {
		response.OkWithDetailed(gin.H{
			"ok": false,
		}, "无sk或获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{
		"ok": true,
	}, "获取成功", c)
}

func (chat *ChatGptApi) DeleteSK(c *gin.Context) {
	err := chatGptService.DeleteSK()
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (chat *ChatGptApi) GetTable(c *gin.Context) {
	var req request.ChatGptRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	sql, results, err := chatGptService.GetTable(req)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithDetailed(gin.H{
			"sql":     sql,
			"results": results,
		}, "生成失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{
		"sql":     sql,
		"results": results,
	}, "ChatGpt生成完成", c)
}

func (chat *ChatGptApi) TestConnect(c *gin.Context) {
	var req sysModel.Datasource
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	names, err := chatGptService.TestConnect(req)
	if err != nil {
		global.GVA_LOG.Error("数据库连接失败!", zap.Error(err))
		response.FailWithDetailed(gin.H{}, "连接失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{
		"names": names,
	}, "连接成功", c)
}

func (chat *ChatGptApi) GetSchema(c *gin.Context) {
	var req request.ChatGptRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	names, err := chatGptService.GetSchema(req)
	if err != nil {
		global.GVA_LOG.Error("获取schema失败!", zap.Error(err))
		response.FailWithDetailed(gin.H{}, "获取schema失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{
		"names": names,
	}, "获取schema成功", c)

}
