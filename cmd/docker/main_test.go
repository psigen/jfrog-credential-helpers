package main

import (
	"os"
	"testing"
)

func TestCommandList(t *testing.T) {
	os.Args = []string{
		"docker-credential-helper",
		"list",
	}
	main()
}
