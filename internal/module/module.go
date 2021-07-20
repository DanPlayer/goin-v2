package module

import (
	"flying-star/internal/config"
	"flying-star/internal/module/casbin"
	"flying-star/internal/module/common"
	"flying-star/internal/module/user"
)

func Registry(config config.Options) common.ModuleOptionList {
	return common.ModuleOptionList{
		user.Init(config),
		casbin.Init(config),
	}
}
