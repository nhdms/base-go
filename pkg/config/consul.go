package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// GetRemoteConfig secret was set at key "consul.secret"
func GetRemoteConfig(endpoint, key string, needDecrypt bool) (value []byte, err error) {
	// init config connection to consul
	config := api.DefaultConfig()
	if endpoint != "" {
		config.Address = endpoint
	}

	// init consul client
	client, err := api.NewClient(config)
	if err != nil {
		return
	}
	kv := client.KV()

	// get key
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return
	}
	if pair == nil {
		err = fmt.Errorf("remote config key is not existed: %v", key)
		return
	}

	value = pair.Value
	if needDecrypt {
		valueBytes, err := hex.DecodeString(string(value))
		if err != nil {
			fmt.Printf("cannot decode remote config: %v", key)
		}
		secret := viper.GetString("consul.secret")
		if len(secret) < 1 {
			return nil, errors.Errorf("Cannot get secret from config!")
		}
		secretKey := []byte(secret)
		value, err = DecryptAES(secretKey, valueBytes)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return
}

func DecryptAES(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
