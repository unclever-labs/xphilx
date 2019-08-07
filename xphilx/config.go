package xphilx

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

	if !strings.HasPrefix(cfg.S3BucketPath, "s3://") {
		err = errors.New("config s3bucketpath does not start with s3://")
		return
	}

	if len(strings.TrimSpace(strings.Replace(cfg.S3BucketPath, "s3://", "", -1))) == 0 {
		err = errors.New("config s3bucketpath does not specify a bucket")
		return
	}

	foundRegion := false
	for _, env := range os.Environ() {
		env := strings.ToLower(strings.Replace(env, "_", "-", -1))
		if strings.Contains(env, "aws-region") || strings.Contains(env, "aws-default-region") {
			foundRegion = true
			break
		}
	}

	if !foundRegion {
		err = errors.New("AWS_REGION || AWS_DEFAULT_REGION not in environment")
		return
	}

	return
}
