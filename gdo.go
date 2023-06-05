package qgomysql

import (
	"strings"
)

func (config *Config) ToString() string {
	return config.query
}

func (config *Config) Prepare(query string) {
	config.query = query
}

func (config *Config) BindParam(key string, value interface{}) {
	config.query = strings.Replace(config.query, ":"+key, escapeValue(value), -1)
}

func (config *Config) Execute() string {
	config.getData()
	return config.result
}
