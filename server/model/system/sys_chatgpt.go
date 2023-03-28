package system

type ChatGpt struct {
	DBName string `json:"dbname,omitempty"`
	Chat   string `json:"chat,omitempty"`
	ChatID string `json:"chatID,omitempty"`
}

type SysChatGptOption struct {
	SK string `json:"sk"`
}

type ChatField struct {
	TABLE_NAME     string
	COLUMN_NAME    string
	COLUMN_COMMENT string
}

type ChatFieldNoTable struct {
	COLUMN_NAME    string
	COLUMN_COMMENT string
}

type Datasource struct {
	Url      string `json:"url"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
