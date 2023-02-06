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
	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments")
		usage()
		os.Exit(1)
	}

	port := os.Args[1]
	bds := os.Args[2]
	baudRate, err := strconv.Atoi(bds)
	if err != nil {
		fmt.Println("Invalid baud rate")
		usage()
		os.Exit(1)
	}
	if baudRate <= 0 {
		fmt.Println("Invalid baud rate")
		usage()
		os.Exit(1)
	}

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

	// Wait for kill
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	card.Stop()
	fmt.Println("Scanning stopped")
}

func cb(info *rfid.CallbackInfo) {
	fmt.Println(info.ID)
}

func usage() {
	fmt.Println("/path/to/bin [portName] [baudRate]")
}