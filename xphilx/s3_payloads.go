package xphilx

import "fmt"

func sendPayloadsToS3(payloadCh chan []byte) (err error) {
	fmt.Println("Consuming payloads")

	for {
		payload := <-payloadCh
		fmt.Println("Payload recieved:", string(payload))
	}

	return
}
