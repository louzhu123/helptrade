package Config

type Config struct {
	ApiKey      string      `json:"ApiKey" ini:"ApiKey" `
	SecretKey   string      `json:"SecretKey" ini:"SecretKey"`
	ConfigMysql ConfigMysql `json:"mysql" ini:"mysql"`
}

type ConfigMysql struct {
	Dsn string `json:"dsn"  ini:"dsn" `
}
