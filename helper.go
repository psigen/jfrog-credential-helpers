package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
)

// ArtifactoryKeychain handles secrets using the Artifactory API for access tokens.
type ArtifactoryKeychain struct{}

var ErrInvalidDomain = errors.New("invalid domain: must contain *.jfrog.io")

func getServerIdFromUrl(serverUrl string) (string, error) {
	fmt.Fprintln(os.Stderr, "serverUrl", serverUrl)

	var serverHostname string
	if strings.HasPrefix(serverUrl, "https://") || strings.HasPrefix(serverUrl, "http://") {
		u, err := url.Parse(serverUrl)
		if err != nil {
			return "", err
		}
		serverHostname = u.Hostname()
	} else {
		serverHostname = serverUrl
	}

	serverId, domain, found := strings.Cut(serverHostname, ".")
	if !found || domain != "jfrog.io" {
		return "", ErrInvalidDomain
	}

	return serverId, nil
}

// Logs into artifactory using an interactive webauth flow.
func (h ArtifactoryKeychain) Login(serverURL string) error {
	fmt.Fprintln(os.Stderr, "Login", serverURL)

	// TODO(PV): Logout before logging in.

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
	fmt.Fprintln(os.Stderr, "Logout", serverURL)

	serverId, err := getServerIdFromUrl(serverURL)
	if errors.Is(err, ErrInvalidDomain) {
		fmt.Println("Skipping logout from invalid domain:", serverURL)
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
	fmt.Fprintln(os.Stderr, "Add", creds)
	return errors.New("adding credentials not supported: use `jf login` CLI")
}

// (UNSUPPORTED) Delete removes credentials from the keychain.
func (h ArtifactoryKeychain) Delete(serverURL string) error {
	fmt.Fprintln(os.Stderr, "Delete", serverURL)
	return errors.New("deleting credentials not supported: use `jf logout` CLI")
}

// Get returns the username and secret to use for a given registry server URL.
func (h ArtifactoryKeychain) Get(serverURL string) (string, string, error) {
	fmt.Fprintln(os.Stderr, "Get", serverURL)

	serverId, err := getServerIdFromUrl(serverURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err", err)
		return "", "", err
	}

	serverDetails, err := config.GetSpecificConfig(serverId, false, false)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err", err)
		return "", "", err
	}

	user := serverDetails.User
	pass := serverDetails.AccessToken
	fmt.Fprintln(os.Stderr, "Succ", user, pass)
	return user, pass, nil
}

// List returns the stored URLs and corresponding usernames.
func (h ArtifactoryKeychain) List() (map[string]string, error) {
	fmt.Fprintln(os.Stderr, "List")

	resp := make(map[string]string)

	serverDetailList, err := config.GetAllServersConfigs()
	if err != nil {
		return nil, err
	}

	for _, serverDetails := range serverDetailList {
		serverHostname := fmt.Sprintf("%s.jfrog.io", serverDetails.ServerId)
		serverUrl := fmt.Sprintf("https://%s/v2", serverHostname)
		resp[serverHostname] = serverDetails.User
		resp[serverUrl] = serverDetails.User
	}

	return resp, nil
}
