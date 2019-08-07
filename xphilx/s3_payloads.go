package xphilx

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func consumePayloads(cfg Config, uploader *s3manager.Uploader, payloadCh chan []byte) (err error) {
	fmt.Println("Consuming payloads")

	cleanupCh := make(chan string)

	var f *os.File

	retry(defaultRetries, defaultPeriod, func() (err error) {
		f, err = ioutil.TempFile("", "xphilx")
		if err != nil {
			fmt.Println("Failed creating tmpfile")
			return
		}

		return
	})

	cnt := 0
	for {
		select {
		case payload := <-payloadCh:
			if payload[len(payload)-1] != '\n' {
				payload = append(payload, byte('\n'))
			}

			retry(defaultRetries, defaultPeriod, func() (err error) {
				if _, err = f.Write(payload); err != nil {
					fmt.Println("Failed writing payload to tmp file")
					return
				}
				return
			})

			cnt++
			if cnt == cfg.LogsPerFile {
				cnt = 0

				go s3Upload(uploader, f.Name(), cfg.S3BucketPath, cleanupCh)

				retry(defaultRetries, defaultPeriod, func() (err error) {
					f, err = ioutil.TempFile("", "xphilx")
					if err != nil {
						fmt.Println("Failed creating tmpfile")
						return
					}

					return
				})
			}
		case fileName := <-cleanupCh:
			fmt.Println("Cleaning up file:", fileName)
			os.Remove(fileName)
		}

	}

	return
}
