package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type ChatGptService struct{}

func (chat *ChatGptService) CreateSK(option system.SysChatGptOption) error {
	_, err := chat.GetSK()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Create(option).Error
		}
		return err
	}
	return errors.New("已经存在sk")
}

func (chat *ChatGptService) GetSK() (option system.SysChatGptOption, err error) {
	err = global.GVA_DB.First(&option).Error
	return
}

func (chat *ChatGptService) DeleteSK() error {
	option, err := chat.GetSK()
	if err != nil {
		return err
	}
	return global.GVA_DB.Delete(option, "sk = ?", option.SK).Error
}

func (chat *ChatGptService) GetTable(req request.ChatGptRequest) (sql2 string, results []map[string]interface{}, err error) {
	if req.DBName == "" {
		return "", nil, errors.New("未选择db")
	}
	//测试连接
	var db *sql.DB
	var jsonBytes []byte
	if req.DbType == "MySql" {
		db, err = utils.TestMysqlConnect(req.UserName, req.Password, req.Url, req.DBName)
		if err != nil {
			return "", nil, err
		}
		//查询mysql库下表名和字断
		jsonBytes, err = utils.GetMysqlDbField(db, req.DBName)
	} else {
		db, err = utils.TestPgSqlConnect(req.UserName, req.Password, req.Url, req.DBName)
		if err != nil {
			return "", nil, err
		}
		//查询pgsql库下表名和字断
		jsonBytes, err = utils.GetMysqlDbField(db, req.Schema)
	}
	defer db.Close()
	if err != nil {
		return "", nil, err
	}
	option, err := chat.GetSK()
	if err != nil {
		return "", nil, err
	}
	client := openai.NewClient(option.SK)
	ctx := context.Background()

	chatReq := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("数据库所有字段用json表示,表名为TABLE_NAME,列名为COLUMN_NAME,列描述为COLUMN_COMMENT,%s,根据语句帮我生成单纯的查询sql,,不要提示语\n+%s", string(jsonBytes), req.Chat),
			},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	resultMap, err := utils.GetResultMap(db, resp.Choices[0].Message.Content)
	//err = global.GVA_DB.Raw(resp.Choices[0].Message.Content).Scan(&results).Error
	return resp.Choices[0].Message.Content, resultMap, err
}

func (chat *ChatGptService) TestConnect(req system.Datasource) (names []string, err error) {
	var db *sql.DB
	var rows *sql.Rows
	if req.DbType == "MySql" {
		db, err = utils.TestMysqlConnect(req.UserName, req.Password, req.Url, "")
		if err != nil {
			return
		}
		rows, err = db.Query("SELECT schema_name FROM information_schema.schemata")

	} else {
		db, err = utils.TestPgSqlConnect(req.UserName, req.Password, req.Url, "")
		if err != nil {
			return
		}
		rows, err = db.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
	}
	defer db.Close()
	if err != nil {
		return
	}
	defer rows.Close()
	var schemaNames = make([]string, 0)
	for rows.Next() {
		var schemaName string
		err = rows.Scan(&schemaName)
		if err != nil {
			panic(err.Error())
		}
		schemaNames = append(schemaNames, schemaName)
	}
	return schemaNames, err
}

func (chat *ChatGptService) GetSchema(req request.ChatGptRequest) (names []string, err error) {
	db, err := utils.TestPgSqlConnect(req.UserName, req.Password, req.Url, req.DBName)
	if err != nil {
		return
	}
	rows, err := db.Query("SELECT schema_name FROM information_schema.schemata")
	defer db.Close()
	if err != nil {
		return
	}
	defer rows.Close()
	var schemaNames = make([]string, 0)
	for rows.Next() {
		var schemaName string
		err = rows.Scan(&schemaName)
		if err != nil {
			panic(err.Error())
		}
		schemaNames = append(schemaNames, schemaName)
	}
	return schemaNames, err
}
