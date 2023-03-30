package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	_ "github.com/lib/pq"
)

func TestMysqlConnect(username, password, url, dbName string) (sqlDb *sql.DB, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, url, dbName))
	if err != nil {
		return
	}
	// 测试连接
	return db, db.Ping()
}

func TestPgSqlConnect(username, password, url, dbName string) (sqlDb *sql.DB, err error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", username, password, url, dbName))
	if err != nil {
		return
	}
	// 测试连接
	return db, db.Ping()
}
func GetMysqlDbField(db *sql.DB, name string) ([]byte, error) {
	sprintf := fmt.Sprintf("SELECT table_name, column_name FROM information_schema.columns WHERE table_schema = '%s'", name)
	rows, _ := db.Query(sprintf)
	defer rows.Close()
	var tablesInfo []system.ChatField
	for rows.Next() {
		var tableName, columnName string
		rows.Scan(&tableName, &columnName)
		tablesInfo = append(tablesInfo, system.ChatField{TABLE_NAME: tableName, COLUMN_NAME: columnName})
	}
	return json.Marshal(tablesInfo)
}

func GetResultMap(db *sql.DB, message string) (results []map[string]interface{}, err error) {
	gptRow, err := db.Query(message)
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
	return resultsMap, err
}
