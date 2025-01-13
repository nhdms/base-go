package config

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type ConfigType string

const (
	ConfigTypeToml ConfigType = "toml"
	ConfigTypeYAML ConfigType = "yaml"
	ConfigTypeJSON ConfigType = "json"
)

func GetConsulAddr() string {
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if len(consulAddr) == 0 {
		consulAddr = "consul:8500"
	}
	return consulAddr
}

func LoadConfigFromFolder(folderPath string, configType ConfigType) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only load files with the specified extension
		if !info.IsDir() && filepath.Ext(info.Name()) == "."+string(configType) {
			v := viper.New()
			v.SetConfigFile(path)
			v.SetConfigType(string(configType))

			if err := v.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading config file %s: %w", path, err)
			}

			if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
				return fmt.Errorf("error merging config file %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error loading config files from folder %s: %w", folderPath, err)
	}
	return nil
}

func LoadConfigFromConsul(keys []string, addr string) error {
	config := api.DefaultConfig()
	if len(addr) == 0 {
		config.Address = GetConsulAddr()
	}

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	kv := client.KV()
	viper.SetConfigType("toml")
	for _, key := range keys {
		pair, _, err := kv.Get(key, nil)
		if err != nil {
			return fmt.Errorf("%v %v", key, err.Error())
		}

		if pair != nil {
			if err := viper.MergeConfig(bytes.NewReader(pair.Value)); err != nil {
				return fmt.Errorf("%v %v", key, err.Error())
			}
		}
		log.Println("loaded key", key)
	}

	return nil
}

func LoadViperToVar(a interface{}, groupKey string) {
	obj := reflect.ValueOf(a).Elem()
	t := reflect.TypeOf(a).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := obj.Field(i)
		field := t.Field(i)
		if f.Kind() != reflect.Struct {
			continue
		}

		tag := field.Tag.Get("json")
		if len(tag) == 0 {
			continue
		}

		if tag == "-" {
			continue
		}

		if !f.CanAddr() {
			continue
		}

		readViperByGroup(f.Addr().Interface(), field.Tag.Get("json"))
	}
	readViperByGroup(a, groupKey)
	//readViperByGroup(&a.Postgres, "postgres")
}

func readViperByGroup[T any](p T, s string) {
	values := viper.GetStringMap(s)
	obj := reflect.ValueOf(p).Elem()
	t := reflect.TypeOf(p).Elem()
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i)
		jsonFileName := fieldName.Tag.Get("json")
		defaultValue := fieldName.Tag.Get("env-default")
		val := obj.Field(i)

		if val.Kind() == reflect.Struct {
			readViperByGroup(val.Addr().Interface(), fmt.Sprintf("%v.%v", s, jsonFileName))
			return
		}

		if v, ok := values[jsonFileName]; ok {
			val.Set(convertToType(v, val.Type()))
			continue
		}

		val.Set(convertToType(defaultValue, val.Type()))
	}
}

// convertToType converts the value to the specified type
func convertToType(value interface{}, targetType reflect.Type) reflect.Value {
	// Get the reflect.Value of the value
	valueReflect := reflect.ValueOf(value)
	if !reflect.TypeOf(value).ConvertibleTo(targetType) {
		if targetType == reflect.TypeOf(time.Duration(0)) {
			durationValue, err := time.ParseDuration(cast.ToString(value))
			if err == nil {
				return reflect.ValueOf(durationValue)
			}
		}
		switch targetType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Convert the string value to int
			intValue, err := strconv.Atoi(valueReflect.String())
			if err != nil {
				fmt.Println(targetType.Name())
				return reflect.Value{}
			}
			// Convert the int value to the target type
			convertedValue := reflect.ValueOf(intValue).Convert(targetType)
			return convertedValue
		}
	}
	// Convert the value to the target type
	convertedValue := valueReflect.Convert(targetType)

	return convertedValue
}

func LoadConfigToVar(cfg interface{}, sub string) error {
	subViper := viper.Sub(sub)
	if subViper == nil {
		return fmt.Errorf("sub-config not found: %s", sub)
	}

	return subViper.Unmarshal(cfg)
}

func ViperGetInt64WithDefault(key string, defaultValue int64) int64 {
	v := viper.GetInt64(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func ViperGetIntWithDefault(key string, defaultValue int) int {
	v := viper.GetInt(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

func ViperGetStringWithDefault(key string, defaultValue string) string {
	v := viper.GetString(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}

func ViperGetDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	v := viper.GetDuration(key)
	if v == 0 {
		return defaultValue
	}
	return v
}

// LoadConfigFromString LoadConfigFromBytes loads configuration from a byte slice into a new Viper instance.
// The configType parameter should specify the format (e.g., "yaml", "json", "toml")
func LoadConfigFromString(configString string, configType string) error {
	viper.SetConfigType(configType)

	// Create a bytes reader
	reader := bytes.NewReader([]byte(configString))

	// Read the config from the byte slice
	if err := viper.ReadConfig(reader); err != nil {
		return err
	}

	return nil
}
