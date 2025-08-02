package common

import (
	"errors"
	"github.com/nhdms/base-go/pkg/utils"
)

type ServiceConfig struct {
	AppType          string   `json:"app_type" mapstructure:"app_type"`
	CmdBinDir        string   `json:"cmd_bin_dir" mapstructure:"cmd_bin_dir"`
	ServiceName      string   `json:"service_name" mapstructure:"service_name"`
	Port             int      `json:"port" mapstructure:"port"`
	ConfigRemoteKeys []string `json:"config_remote_keys" mapstructure:"config_remote_keys"`
	ConsulAddrs      []string `json:"consul_addrs" mapstructure:"consul_addrs"`
}

const (
	ServiceTypeAPI        = "api"
	ServiceTypeConsumer   = "consumer"
	ServiceTypeService    = "service"
	ServiceTypeSchedulers = "schedulers"
)

var ValidServiceTypes = map[string]bool{
	ServiceTypeAPI:        true,
	ServiceTypeConsumer:   true,
	ServiceTypeService:    true,
	ServiceTypeSchedulers: true,
}

func (c *ServiceConfig) Validate() error {
	if len(c.ServiceName) == 0 {
		return errors.New("service_name is required")
	}

	if len(c.AppType) == 0 {
		return errors.New("app_type is required")
	}

	if _, ok := ValidServiceTypes[c.AppType]; !ok {
		return errors.New("invalid app_type")
	}

	return nil
}

func (c *ServiceConfig) ToString() string {
	return utils.ToJSONString(c)
}
