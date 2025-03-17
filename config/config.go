package Config

type Config struct {
	ApiKey       string      `json:"ApiKey" ini:"ApiKey" `
	SecretKey    string      `json:"SecretKey" ini:"SecretKey"`
	MiaoNoticeId string      `json:"MiaoNoticeId" ini:"MiaoNoticeId"`
	ConfigMysql  ConfigMysql `json:"mysql" ini:"mysql"`
}

type ConfigMysql struct {
	Dsn string `json:"dsn"  ini:"dsn" `
}
