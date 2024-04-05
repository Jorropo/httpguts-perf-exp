package p

import (
	"strings"

	"golang.org/x/net/http/httpguts"
)

func IsCookieNameValid(raw string) bool {
	return IsToken(raw)
}

// implementation copied from https://github.com/golang/go/blob/2e064cf14441460290fd25d9d61f02a9d0bae671/src/net/http/cookie.go#L463
func IsCookieNameValidStd(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}

// implementation copied from https://github.com/golang/go/blob/2e064cf14441460290fd25d9d61f02a9d0bae671/src/net/http/http.go#L61
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}
