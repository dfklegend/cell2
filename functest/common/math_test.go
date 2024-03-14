package common

import (
	"log"
	"testing"
)

func Test_Math(t *testing.T) {
	f := 1.1
	log.Printf("%v\n", int(f))
	f = 1.9
	log.Printf("%v\n", int(f))
	f = 1111111111111111110000.9
	log.Printf("%v\n", int(f))
	log.Printf("%v\n", int64(f))
}
