package xphilx

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3Upload(uploader *s3manager.Uploader, fileName, s3BucketPath string, cleanupCh chan string) {
	retry(defaultRetries, defaultPeriod, func() (err error) {
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Failed opening file")
			return
		}

		bucket, key := getBucketKey(s3BucketPath, false)
		uploadInput := s3manager.UploadInput{
			Body:   f,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}

		if _, err = uploader.Upload(&uploadInput); err != nil {
			fmt.Println("Failed uploading to s3")
			return
		}

		cleanupCh <- fileName

		return
	})
}

func getBucketKey(s3BucketPath string, test bool) (bucket, key string) {
	s3BucketPath = strings.Replace(s3BucketPath, "s3://", "", -1)
	parts := strings.Split(s3BucketPath, "/")

	bucket = strings.TrimSuffix(parts[0], "/")
	keyPrefix := strings.TrimSuffix(strings.Join(parts[1:], "/"), "/")

	now := time.Now().UTC()
	dateStr := now.Format("2006/01/02")
	timeStr := now.Format("payload-03-04-05.log")

	key = strings.TrimPrefix(fmt.Sprintf("%s/%s/%s", keyPrefix, dateStr, timeStr), "/")
	if test {
		key = strings.TrimPrefix(fmt.Sprintf("%s/test.log", keyPrefix), "/")
	}
	return
}

func testUpload(cfg Config, uploader *s3manager.Uploader) (err error) {
	buf := bytes.NewBufferString(`{"msg":"test string"}`)

	bucket, key := getBucketKey(cfg.S3BucketPath, true)
	in := &s3manager.UploadInput{
		Body:   buf,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	fmt.Printf("Testing upload of %s/%s\n", bucket, key)
	if _, err = uploader.Upload(in); err != nil {
		fmt.Println("Failed test upload")
		return
	}

	return
}
