# Olimex RFID/NFC readers Golang lib

Allow to connect to serial port exposed by olimex rfids/nfc readers and scanning card ids.

## Usage

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	rfid "github.com/angelodlfrtr/olimex-rfid-serial-go"
)

func main() {
	port := "/dev/ttyWHATYW"
	baudRate := 9600
	card, err := rfid.New(port, baudRate, cb)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := card.Scan(); err != nil {
			if !strings.Contains(err.Error(), "Port has been closed") {
				panic(err)
			}
		}
	}()

	// Wait for ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

  // Stop scan
	card.Stop()
	fmt.Println("Scan stopped")
}

func cb(info *rfid.CallbackInfo) {
	fmt.Println(info.ID)
}
```