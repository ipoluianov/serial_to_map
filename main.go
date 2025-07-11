package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ipoluianov/gomisc/logger"
	"github.com/tarm/serial"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <serial-port>")
		return
	}

	serialPort := os.Args[1]

	config := &serial.Config{
		Name: serialPort,
		Baud: 9600,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		logger.Println("open port error:", err)
	}
	defer port.Close()

	input := make([]byte, 128)
	for {
		buffer := make([]byte, 128)
		n, err := port.Read(buffer)
		if err != nil {
			logger.Println("read error:", err)
			break
		}
		if n > 0 {
			input = append(input, buffer[:n]...)
		}
		indexOf := strings.Index(string(input), "\n")
		if indexOf == -1 {
			continue
		}
		frame := input[:indexOf+1]
		line := strings.TrimFunc(string(frame), func(r rune) bool {
			return r == '\r' || r == '\n' || r == '\x00'
		})
		input = input[indexOf+1:]
		logger.Println("line:", line)
		tempValueParts := strings.Split(line, "=")
		if len(tempValueParts) != 2 {
			logger.Println("invalid line:", line)
			continue
		}
		tempValue := tempValueParts[1]

		_, err = http.Get("http://map.u00.io/set/t/" + tempValue)
		if err != nil {
			logger.Println("http error:", err)
			continue
		}

		time.Sleep(500 * time.Millisecond)
	}
}
