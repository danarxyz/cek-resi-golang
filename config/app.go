package config

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

func Boot() {}

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		"name":            config.Env("APP_NAME", "Goravel"),
		"env":             config.Env("APP_ENV", "production"),
		"debug":           config.Env("APP_DEBUG", false),
		"timezone":        carbon.UTC,
		"locale":          "en",
		"fallback_locale": "en",
		"key":             config.Env("APP_KEY", ""),
	})
}
