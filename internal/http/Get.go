package http

import (
	"errors"
	"io"
	"net/http"
)

func Get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("GET request error: " + err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("response read error: " + err.Error())
	}

	return body, nil
}
