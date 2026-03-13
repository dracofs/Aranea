package utils

import (
	"net/url"
)

func normalize(base string, link string) (string, error) {
	baseURL, err := url.Parse(link)

	if err != nil {
		return "", err
	}

	refURL, err := url.Parse(link)

	if err != nil {
		return "", err
	}

	resolved := baseURL.ResolveReference(refURL).String()
	return resolved, nil
}
