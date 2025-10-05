package utils

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func IsValidURL(str string) bool {
	if str == "" {
		return false
	}

	// contains scheme
	if !(len(str) > 7 && (str[:7] == "http://" || str[:8] == "https://")) {
		return false
	}

	_, err := url.ParseRequestURI(str)
	return err == nil
}

func IsUUID(str string) bool {
	if len(str) == 36 {
		if len(strings.ReplaceAll(str, "-", "")) == 32 {
			return true
		}
	}

	return false
}

func GetDataFromBase64(encoded string) (*[]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetDataFromURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch file: " + resp.Status)
	}

	return resp.Body, nil
}

func BoolPtr(v bool) *bool {
	return &v
}

func IntPtr(v int) *int {
	return &v
}

func StringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
