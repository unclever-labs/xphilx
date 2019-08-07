package xphilx

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// Config is configuration for the worker
type Config struct {
	Interface    string `validate:"required"`
	LogsPerFile  int    `validate:"required"`
	Port         string `validate:"required"`
	S3BucketPath string `validate:"required"`
	SnapLength   int32  `validate:"required"`
}

func validateConfig(cfg Config) (err error) {
	if err = validator.New().Struct(cfg); err != nil {
		fmt.Println("Config is invalid and probably missing fields")
		return
	}

	return
}
