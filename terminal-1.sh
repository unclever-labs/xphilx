#!/usr/bin/env bash

# sudo AWS_REGION=us-east-1 AWS_PROFILE=update go run main.go -b s3://update-bucket-name-here -i lo0 -l 3
sudo AWS_REGION=us-west-2 AWS_PROFILE=personal go run main.go -b s3://rms1000watt-test-bucket -i lo0 -l 3
