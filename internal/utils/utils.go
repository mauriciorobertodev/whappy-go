package utils

import (
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

func Uint32Ptr(v uint32) *uint32 {
	return &v
}
