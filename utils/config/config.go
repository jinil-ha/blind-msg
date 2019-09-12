package config

import (
	"crypto/tls"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper
}

var globalConfig *Config
var environment string

func init() {
	configName := os.Getenv("CAR_MSG_CONFIG")
	if len(configName) == 0 {
		configName = "server"
	}

	var err error
	globalConfig, err = ParseFile(configName)

	if err != nil {
		panic(err)
	}

	environment = globalConfig.GetString("environment")

	// TODO: check required config
}

func ParseFile(fn string) (*Config, error) {
	v := viper.New()

	v.AddConfigPath("/etc/carmsg/")
	//v.AddConfigPath("$HOME/.e2e-server")
	v.AddConfigPath(".")

	v.SetConfigName(fn)
	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return &Config{
		v: v,
	}, nil
}

func Get(key string) interface{} {
	return globalConfig.Get(key)
}

func GetString(key string) string {
	return globalConfig.GetString(key)
}

func GetStringList(key string) []string {
	return globalConfig.GetStringList(key)
}

func GetInt(key string) int {
	return globalConfig.GetInt(key)
}

func GetBool(key string) bool {
	return globalConfig.GetBool(key)
}

func GetTLSCertificate(certKey, keyKey string) (tls.Certificate, error) {
	return globalConfig.GetTLSCertificate(certKey, keyKey)
}

func Unmarshal(key string, rawVal interface{}) error {
	return globalConfig.Unmarshal(key, rawVal)
}

func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) GetStringList(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Config) GetTLSCertificate(certKey, keyKey string) (certificate tls.Certificate, err error) {
	cert := c.GetString(certKey)
	key := c.GetString(keyKey)
	if strings.HasPrefix(cert, "-----") {
		certificate, err = tls.X509KeyPair([]byte(cert), []byte(key))
	} else {
		certificate, err = tls.LoadX509KeyPair(cert, key)
	}
	return
}

func (c *Config) Unmarshal(key string, rawVal interface{}) error {
	return c.v.UnmarshalKey(key, rawVal)
}

//-------------------------- API for test

func InitForUnitTest() {
	globalConfig = &Config{
		v: viper.New(),
	}
}

func (c *Config) SetStringForUnitTest(key string, value string) {
	c.v.Set(key, value)
}

func SetStringForUnitTest(key string, value string) {
	globalConfig.SetStringForUnitTest(key, value)
}

func IsProduction() bool {
	return environment == "production"
}

func IsStage() bool {
	return environment == "stage"
}

func IsDev() bool {
	return environment == "dev" || environment == "beta" || environment == "alpha"
}

func IsNotProduction() bool {
	return !IsProduction()
}

func Environment() string {
	return environment
}
