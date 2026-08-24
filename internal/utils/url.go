package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func Normalize(base string, link string) (string, error) {
	link = strings.TrimSpace(link)
	if link == "" {
		return "", fmt.Errorf("empty link")
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	refURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	resolved := baseURL.ResolveReference(refURL)
	resolved.Fragment = ""

	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme %q", resolved.Scheme)
	}
	if resolved.Host == "" {
		return "", fmt.Errorf("missing host")
	}

	return resolved.String(), nil
}
