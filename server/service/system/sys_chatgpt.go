package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
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
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", req.UserName, req.Password, req.Url, req.DBName))
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	sprintf := fmt.Sprintf("SELECT table_name, column_name FROM information_schema.columns WHERE table_schema = '%s'", req.DBName)
	rows, _ := db.Query(sprintf)
	defer rows.Close()
	var tablesInfo []system.ChatField
	for rows.Next() {
		var tableName, columnName string
		rows.Scan(&tableName, &columnName)
		tablesInfo = append(tablesInfo, system.ChatField{TABLE_NAME: tableName, COLUMN_NAME: columnName})
	}
	jsonBytes, err := json.Marshal(tablesInfo)
	if err != nil {
		return "", nil, err
	}
	//var tablesInfo []system.ChatField
	//global.GVA_DB.Table("information_schema.columns").Where("TABLE_SCHEMA = ?", req.DBName).Scan(&tablesInfo)
	//b, err := json.Marshal(tablesInfo)
	//if err != nil {
	//	return
	//}
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
	gptRow, err := db.Query(resp.Choices[0].Message.Content)
	if err != nil {
		return
	}
	// 获取列名
	columns, err := gptRow.Columns()
	if err != nil {
		return
	}
	resultsMap := make([]map[string]interface{}, 0)

	// 遍历查询结果
	for gptRow.Next() {
		// 创建一个map，用于存储每行数据的值
		resultMap := make(map[string]interface{})

		// 创建一个切片，用于存储每行数据的值
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// 将查询结果扫描到values切片中
		err = gptRow.Scan(values...)
		if err != nil {
			return
		}

		// 将每行数据的值存储到map中
		for i, column := range columns {
			value := *(values[i].(*interface{}))
			// 将值转换为[]byte类型，并将其转换为字符串
			var strValue string
			switch v := value.(type) {
			case []byte:
				strValue = string(v)
			default:
				strValue = fmt.Sprint(v)
			}

			resultMap[column] = strValue
		}

		// 将map添加到切片中
		resultsMap = append(resultsMap, resultMap)
	}

	defer gptRow.Close()
	//err = global.GVA_DB.Raw(resp.Choices[0].Message.Content).Scan(&results).Error
	return resp.Choices[0].Message.Content, resultsMap, err
}

func (chat *ChatGptService) TestConnect(req system.Datasource) (names []string, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", req.UserName, req.Password, req.Url))
	if err != nil {
		return
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		return
	}
	rows, err := db.Query("SELECT schema_name FROM information_schema.schemata")
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
