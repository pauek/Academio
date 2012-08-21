package data

import (
	"log"
	"fmt"
	"os"
)

func NewUUID() string {
	f, err := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Cannot open /dev/urandom: %s", err)
	}
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}