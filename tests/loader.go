package tests

import (
	"github.com/nhdms/base-go/pkg/config"
	"path"
	"path/filepath"
	"runtime"
)

func LoadTestConfig() error {
	x := RootDir()
	return config.LoadConfigFromFolder(x+"/tests/config", config.ConfigTypeToml)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
