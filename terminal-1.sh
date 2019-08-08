#!/usr/bin/env bash

sudo AWS_REGION=us-east-1 AWS_PROFILE=update go run main.go -b s3://update-bucket-name-here -i lo0 -l 3
