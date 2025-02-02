package main

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/psigen/jfrog-credential-helpers/internal/bazel"
	"github.com/psigen/jfrog-credential-helpers/internal/jfrog"
)

type ArtifactoryCredentialProvider struct{}

// Get the credentials for a specific JFrog URI.
func (h ArtifactoryCredentialProvider) Get(request bazel.GetCredentialsRequest) (*bazel.GetCredentialsReponse, error) {
	serverURL := request.Uri

	serverHostname, err := jfrog.GetHostnameFromURLorHost(serverURL)
	if err != nil {
		log.Fatalln("unable to find hostname in", serverURL)
		return nil, err
	}

	serverDetailList, err := config.GetAllServersConfigs()
	if err != nil {
		log.Fatalln("unable to retrieve JFrog server configs:", err)
		return nil, err
	}

	for _, serverDetails := range serverDetailList {
		u, err := url.Parse(serverDetails.Url)
		if err != nil {
			log.Println("Skipping", serverDetails.Url, ": unable to parse as URL")
			continue
		}

		if u.Hostname() == serverHostname {
			var response bazel.GetCredentialsReponse

			user := serverDetails.User
			pass := serverDetails.AccessToken

			authStr := b64.URLEncoding.EncodeToString([]byte(
				fmt.Sprintf("%s:%s", user, pass),
			))

			response.Headers = map[string][]string{
				"Authorization": {"Basic " + authStr},
			}

			return &response, nil
		}
	}

	log.Fatalln("not logged into", serverURL, "from JFrog CLI")
	return nil, errors.ErrUnsupported
}
