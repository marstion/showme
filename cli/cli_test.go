package cli

import (
	"log"
	"testing"
)

func TestCli(t *testing.T) {
	outStr, err := Cli("curl", []string{"-v", "httpbin.org/ip"})
	if err != nil {
		log.Fatal(err)

	}
	log.Printf("outStr: %#v\n", outStr)
}
