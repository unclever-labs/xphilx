package xphilx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
	payloadCh      chan []byte
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading stream: %v %v: %s", h.net, h.transport, err)
			return
		}

		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Failed reading all bytes from request body:", err)
			return
		}

		h.payloadCh <- payload

		fmt.Println("Payload length:", len(payload))
	}
}
