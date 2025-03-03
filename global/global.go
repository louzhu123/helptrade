package global

import (
	Config "helptrade/config"

	"gorm.io/gorm"
)

var Cfg *Config.Config
var DB *gorm.DB
