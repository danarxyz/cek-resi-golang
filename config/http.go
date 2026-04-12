package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		"url":             config.Env("APP_URL", "http://localhost"),
		"host":            config.Env("APP_HOST", "127.0.0.1"),
		"port":            config.Env("APP_PORT", "3000"),
		"request_timeout": 30,
	})
}
