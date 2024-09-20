package config

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	DEFAULT_CONFIG_PATH = "conf"
	DEFAULT_CONFIG_NAME = "zero-framework"
	DEFAULT_CONFIG_TYPE = "yaml"
)

var (
	_servAbsPath string
	_aViper      *viper.Viper
)

func cfgString(aViper *viper.Viper, cfgName string) string {
	envName := strings.ToUpper(strings.Replace(cfgName, ".", "_", -1))
	envValue := os.Getenv(envName)
	if len(envValue) > 0 {
		return envValue
	}
	return aViper.GetString(cfgName)
}

func cfgInt(aViper *viper.Viper, cfgName string) int {
	envName := strings.ToUpper(strings.Replace(cfgName, ".", "_", -1))
	envValue := os.Getenv(envName)
	if len(envValue) > 0 {
		v, err := strconv.Atoi(envValue)
		if err != nil {
			panic(err)
		}
		return v
	}
	return aViper.GetInt(cfgName)
}

func cfgStringSlice(aViper *viper.Viper, cfgName string) []string {
	envName := strings.ToUpper(strings.Replace(cfgName, ".", "_", -1))
	envValue := os.Getenv(envName)
	if len(envValue) > 0 {
		return strings.Split(envValue, ",")
	}
	return aViper.GetStringSlice(cfgName)
}

func cfgStringMapString(aViper *viper.Viper, cfgName string) map[string]string {
	return aViper.GetStringMapString(cfgName)
}

func cfgStringMapStringSlice(aViper *viper.Viper, cfgName string) map[string][]string {
	return aViper.GetStringMapStringSlice(cfgName)
}

func NewConfigs(sysAbsPath string) {
	if len(_servAbsPath) > 0 {
		return
	}
	_servAbsPath = sysAbsPath
	_aViper = viper.New()

	_aViper.AddConfigPath(path.Join(_servAbsPath, DEFAULT_CONFIG_PATH))

	_aViper.SetConfigName(DEFAULT_CONFIG_NAME)
	_aViper.SetConfigType(DEFAULT_CONFIG_TYPE)

	if err := _aViper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func ServerAbsPath() string {
	return _servAbsPath
}

func StringValue(cfgName string) string {
	return cfgString(_aViper, cfgName)
}

func IntValue(cfgName string) int {
	return cfgInt(_aViper, cfgName)
}

func SliceStringValue(cfgName string) []string {
	return cfgStringSlice(_aViper, cfgName)
}

func StringMapString(cfgName string) map[string]string {
	return cfgStringMapString(_aViper, cfgName)
}

func StringMapStringSlice(cfgName string) map[string][]string {
	return cfgStringMapStringSlice(_aViper, cfgName)
}
