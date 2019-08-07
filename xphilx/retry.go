package xphilx

import (
	"fmt"
	"os"
	"time"
)

func retry(retries int, period time.Duration, fn func() (err error)) {
	err := fn()
	if err == nil {
		return
	}

	fmt.Println("Error:", err)

	retries--

	if retries == -1 {
		fmt.Println("No more retries left")
		os.Exit(1)
	}

	fmt.Println("Sleeping", period)
	time.Sleep(period)
	retry(retries, period, fn)
}
