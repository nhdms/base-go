package tests

import (
	"gitlab.com/a7923/athena-go/pkg/config"
	"gitlab.com/a7923/athena-go/pkg/utils"
)

func LoadTestConfig() error {
	x := utils.RootDir()
	return config.LoadConfigFromFolder(x+"/tests/config", config.ConfigTypeToml)
}
