package helpers

import (
	"errors"
	"net/url"
	"strings"
)

var ErrInvalidDomain = errors.New("invalid domain: must contain *.jfrog.io")

// Returns the hostname of either a URL or a hostname string.
func GetHostnameFromURLorHost(serverUrl string) (string, error) {
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
func GetServerIdFromUrl(serverURL string) (string, error) {
	serverHostname, err := GetHostnameFromURLorHost(serverURL)
	if err != nil {
		return "", err
	}

	serverId, domain, found := strings.Cut(serverHostname, ".")
	if !found || domain != "jfrog.io" {
		return "", ErrInvalidDomain
	}

	return serverId, nil
}
