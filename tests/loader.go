package tests

import (
	"github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/utils"
)

func LoadTestConfig() error {
	x := utils.RootDir()
	return config.LoadConfigFromFolder(x+"/tests/config", config.ConfigTypeToml)
}
