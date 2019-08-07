package xphilx

import (
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

		bucket, key := getBucketKey(s3BucketPath)
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

func getBucketKey(s3BucketPath string) (bucket, key string) {
	s3BucketPath = strings.Replace(s3BucketPath, "s3://", "", -1)
	parts := strings.Split(s3BucketPath, "/")

	bucket = parts[0]
	keyPrefix := strings.Join(parts[1:], "/")

	now := time.Now().UTC()
	timeStr := now.Format("payload-03-04-05.log")

	year, monthStr, day := now.Date()
	month := int(monthStr)

	key = fmt.Sprintf("%s/%d/%d/%d/%s", keyPrefix, year, month, day, timeStr)
	return
}
