package global

import (
	Config "main/config"

	"gorm.io/gorm"
)

var Cfg *Config.Config
var DB *gorm.DB
