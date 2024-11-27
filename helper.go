package main

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
)

// ArtifactoryKeychain handles secrets using the Artifactory API for access tokens.
type ArtifactoryKeychain struct{}

var ErrInvalidDomain = errors.New("invalid domain: must contain *.jfrog.io")

// Returns the hostname of either a URL or a hostname string.
func getHostnameFromURLorHost(serverUrl string) (string, error) {
	if strings.HasPrefix(serverUrl, "https://") || strings.HasPrefix(serverUrl, "http://") {
		u, err := url.Parse(serverUrl)
		if err != nil {
			return "", err
		}
		return u.Hostname(), nil
	} else {
		return serverUrl, nil
	}

}

// Gets the default server ID used by JFrog for an artifactory URL.
func getServerIdFromUrl(serverURL string) (string, error) {
	serverHostname, err := getHostnameFromURLorHost(serverURL)
	if err != nil {
		return "", err
	}

	serverId, domain, found := strings.Cut(serverHostname, ".")
	if !found || domain != "jfrog.io" {
		return "", ErrInvalidDomain
	}

	return serverId, nil
}

// Logs into artifactory using an interactive webauth flow.
func (h ArtifactoryKeychain) Login(serverURL string) error {
	// TODO(PV): Logout before logging in?

	serverDetails := config.ServerDetails{Url: serverURL}
	serverId, err := getServerIdFromUrl(serverURL)
	if err != nil {
		return err
	}

	err = commands.NewConfigCommand(commands.AddOrEdit, serverId).
		SetInteractive(true).SetUseWebLogin(true).
		SetDetails(&serverDetails).Run()
	if err != nil {
		return err
	}

	return nil
}

// Logs out of artifactory by clearing the creds used by `jf` CLI tool.
func (h ArtifactoryKeychain) Logout(serverURL string) error {
	serverId, err := getServerIdFromUrl(serverURL)
	if errors.Is(err, ErrInvalidDomain) {
		log.Println("Skipping logout from invalid domain:", serverURL)
		return nil
	} else if err != nil {
		return err
	}

	err = commands.NewConfigCommand(commands.Delete, serverId).Run()
	if err != nil {
		return err
	}

	return nil
}

// (UNSUPPORTED) Add adds new credentials to the keychain.
func (h ArtifactoryKeychain) Add(creds *credentials.Credentials) error {
	log.Fatalln("adding credentials is not supported: use `jf login` CLI")
	return errors.ErrUnsupported
}

// (UNSUPPORTED) Delete removes credentials from the keychain.
func (h ArtifactoryKeychain) Delete(serverURL string) error {
	log.Fatalln("deleting credentials is not supported: use `jf` CLI")
	return errors.ErrUnsupported
}

// Get returns the username and secret to use for a given registry server URL.
func (h ArtifactoryKeychain) Get(serverURL string) (string, string, error) {
	serverHostname, err := getHostnameFromURLorHost(serverURL)
	if err != nil {
		log.Fatalln("unable to find hostname in", serverURL)
	}

	serverDetailList, err := config.GetAllServersConfigs()
	if err != nil {
		log.Fatalln("unable to retrieve JFrog server configs:", err)
		return "", "", err
	}

	for _, serverDetails := range serverDetailList {
		u, err := url.Parse(serverDetails.Url)
		if err != nil {
			log.Println("Skipping", serverDetails.Url, ": unable to parse as URL")
			continue
		}

		if u.Hostname() == serverHostname {
			user := serverDetails.User
			pass := serverDetails.AccessToken
			return user, pass, nil
		}
	}

	log.Fatalln("not logged into", serverURL, "from JFrog CLI")
	return "", "", errors.ErrUnsupported
}

// List returns the stored URLs and corresponding usernames.
func (h ArtifactoryKeychain) List() (map[string]string, error) {
	resp := make(map[string]string)

	serverDetailList, err := config.GetAllServersConfigs()
	if err != nil {
		log.Println("unable to retrieve JFrog server configs:", err)
		return nil, err
	}

	for _, serverDetails := range serverDetailList {
		u, err := url.Parse(serverDetails.Url)
		if err != nil {
			log.Println("Skipping", serverDetails.Url, ": unable to parse as URL")
			continue
		}

		resp[u.Hostname()] = serverDetails.User
	}

	return resp, nil
}
