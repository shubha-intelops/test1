package client

import (
	"io"
	"net/http"
)

func ExecuteNode07() ([]byte, error) {
	response, err := http.Get("login:8087/ping")

	if err != nil {
		return nil, err
	}

	return io.ReadAll(response.Body)
}
