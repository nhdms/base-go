package app

import (
	"github.com/micro/plugins/v5/registry/consul"
	"github.com/nhdms/base-go/pkg/common"
	"github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go-micro.dev/v5/registry"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

const (
	GRPCPrefix = "grpc"
	APIPrefix  = "api"
)

func GetGRPCServiceName(svcName string) string {
	return GRPCPrefix + "." + svcName
}

func GetAPIName(svcName string) string {
	return APIPrefix + "." + svcName
}

var defaultRemoteConfigKeys = map[string][]string{
	common.ServiceTypeAPI:      {"admin/conf.toml", "admin/api.toml"},
	common.ServiceTypeConsumer: {"admin/conf.toml", "admin/consumer.toml"},
	common.ServiceTypeService:  {"admin/conf.toml", "admin/service.toml"},
}

var GlobalServiceConfig *common.ServiceConfig

func init() {
	GlobalServiceConfig = LoadInitConfig()
	logger.InitLogger()
	logger.DefaultLogger.Infow("logger initialized successfully")
}

func LoadInitConfig() *common.ServiceConfig {
	conf := &common.ServiceConfig{}

	ignoreLoadConfig := cast.ToBool(os.Getenv("IGNORE_LOAD_CONFIG"))
	if ignoreLoadConfig {
		return &common.ServiceConfig{}
	}
	// Specify the path and file name
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = "./config"
	}

	// Check if cicd.json exists
	cicdPath := filepath.Join(configPath, "cicd.json")
	if _, err := os.Stat(cicdPath); os.IsNotExist(err) {
		log.Printf("%s not found, loading config from environment variables", cicdPath)

		// Configure Viper for environment variables
		v := viper.New()   // Create a new Viper instance to avoid conflicts
		v.SetEnvPrefix("") // No prefix for env vars
		v.AutomaticEnv()

		// Map struct fields to environment variables
		if err := bindEnvs(v, conf); err != nil {
			log.Fatalf("Error binding environment variables: %v", err)
		}

		// Unmarshal environment variables into the config struct
		if err := v.Unmarshal(conf); err != nil {
			log.Fatalf("Error unmarshaling environment variables: %v", err)
		}

	} else {
		// Load config from cicd.json file
		log.Printf("Loading config from %s", cicdPath)

		v := viper.New() // Create a new Viper instance
		v.SetConfigType("json")
		v.SetConfigFile(cicdPath) // Set the specific config file path

		// Read the config file
		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		// Unmarshal the config into our struct
		if err := v.Unmarshal(conf); err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}
	}

	// Load remote config if service type is specified
	err := conf.Validate()
	if err != nil {
		log.Println("Seems like you are missing config path to cicd.json file\n Please run with CONFIG_PATH=./config go run main.go \n" +
			"or change working directory to folder that contains main.go")
		log.Fatalf("Invalid service config: %v", err)
	}

	log.Println("Config loaded ", conf.ToString())
	serviceType := conf.AppType
	remoteConfigKeys := append(defaultRemoteConfigKeys[serviceType], conf.ConfigRemoteKeys...)
	log.Printf("Remote config keys: %v", remoteConfigKeys)

	err = config.LoadConfigFromConsul(remoteConfigKeys, "")
	if err != nil {
		log.Fatalf("Load config from Consul failed: %v", err)
	}
	log.Printf("Config loaded successfully. Test value(viper.GetInt('test.x')) %v", viper.GetInt("test.x"))

	return conf
}

func GetRegistry() registry.Registry {
	consulAddr := config.GetConsulAddr()
	return consul.NewRegistry(registry.Addrs(consulAddr))
}

// bindEnvs automatically binds environment variables based on struct tags
func bindEnvs(v *viper.Viper, iface interface{}) error {
	t := reflect.TypeOf(iface).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("mapstructure"); tag != "" {
			err := v.BindEnv(tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
