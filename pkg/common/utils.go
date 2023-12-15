package common

import (
	"net/url"
	"strings"
)

const (
	NoValue float64 = -1
)

func ResolveAbsoluteUriReference(baseUrl *url.URL, paths ...*url.URL) *url.URL {
	var relativePath string
	for _, path := range paths {
		relativePath = strings.Join([]string{relativePath, path.String()}, "")
	}

	absoluteUrl := baseUrl.ResolveReference(&url.URL{Path: relativePath})
	return absoluteUrl
}
