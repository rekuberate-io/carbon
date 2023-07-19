package providers

import (
	"net/url"
	"strings"
)

func ResolveAbsoluteUriReference(baseUrl *url.URL, paths ...*url.URL) *url.URL {
	var relativePath string
	for _, path := range paths {
		relativePath = strings.Join([]string{relativePath, path.String()}, "")
	}

	absoluteUrl := baseUrl.ResolveReference(&url.URL{Path: relativePath})
	return absoluteUrl
}
