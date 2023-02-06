package rfid

import (
	"bufio"
	"strings"

	"go.bug.st/serial"
)

// Callback for receiving card infos
type Callback func(*CallbackInfo)

// CallbackInfo contain tag info
type CallbackInfo struct {
	ID string
}

// Card define struct for accessing one of the olimex rfid reader in uart mode
type Card struct {
	port serial.Port
	cb   Callback
}

// New card with given port, baudrate and callback
func New(portName string, baudrate int, cb Callback) (*Card, error) {
	mode := &serial.Mode{
		BaudRate: baudrate,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	card := &Card{port, cb}

	return card, nil
}

// Start scanning
func (card *Card) Scan() error {
	scanner := bufio.NewScanner(card.port)

	// Default scanner split is line detection, so keep the logic
	for scanner.Scan() {
		text := scanner.Text()

		// Each tag ids detection begin with a "-" (minus)
		if strings.HasPrefix(text, "-") {
			cardID := strings.TrimPrefix(text, "-")
			info := &CallbackInfo{ID: cardID}
			card.cb(info)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Stop scanning
func (card *Card) Stop() error {
	return card.port.Close()
}
