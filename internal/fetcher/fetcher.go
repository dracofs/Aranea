package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxBodyBytes = 10 << 20 // 10 MiB

var client = &http.Client{
	Timeout: 15 * time.Second,
}

func Fetch(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Aranea/0.1 (educational crawler)")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 1<<20))
		return "", fmt.Errorf("unexpected status %d", res.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, maxBodyBytes+1))
	if err != nil {
		return "", err
	}
	if len(body) > maxBodyBytes {
		return "", fmt.Errorf("response body exceeds %d bytes", maxBodyBytes)
	}

	return string(body), nil
}
