package global

import (
	"github.com/lgc202/mall-micro/service/user/rpc/config"
	"gorm.io/gorm"
)

var (
	Db           *gorm.DB
	ServerConfig = &config.ServerConfig{}
	NacosConfig  = &config.NacosConfig{}
)
