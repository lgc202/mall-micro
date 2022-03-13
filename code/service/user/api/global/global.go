package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/lgc202/mall-micro/service/user/api/config"
	"github.com/lgc202/mall-micro/service/user/api/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  = &config.ServerConfig{}
	NacosConfig   = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
