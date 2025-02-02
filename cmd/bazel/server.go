package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// See: https://github.com/EngFlow/credential-helper-spec/blob/main/schemas/get-credentials-request.schema.json
type GetCredentialsRequest struct {
	Uri string `json:"uri"`
}

// See: https://github.com/EngFlow/credential-helper-spec/blob/main/schemas/get-credentials-response.schema.json
type GetCredentialsReponse struct {
	Expires string              `json:"expires,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type BazelCredentialProvider interface {
	Get(request GetCredentialsRequest) (*GetCredentialsReponse, error)
}

func Serve(provider BazelCredentialProvider) {
	// See: https://github.com/EngFlow/credential-helper-spec/blob/main/spec.md
	// The only officially-supported command is 'get'.
	switch os.Args[1] {
	case "get":
		var request GetCredentialsRequest

		err := json.NewDecoder(os.Stdin).Decode(&request)
		if err != nil {
			log.Fatalln("unable to decode 'get' request:", err)
		}

		res, err := provider.Get(request)
		if err != nil {
			log.Fatalln("unable to process 'get' request:", request, err)
		}

		enc := json.NewEncoder(os.Stdout)
		err = enc.Encode(res)
		if err != nil {
			log.Fatalln("unable to encode 'get' response:", err)
		}

	default:
		fmt.Println("expected <get> subcommand")
		os.Exit(1)
	}
}
