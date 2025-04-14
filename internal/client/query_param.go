package client

import (
	"fmt"
	netUrl "net/url"
	"strings"
)

// QueryParam represents a query parameter in the URL of an HTTP request.
type QueryParam struct {
	Data // Embeds the Data struct to represent the key-value pair for the query parameter.
}

// setQueryParams appends query parameters to the given URL.
// It returns the URL with the query parameters attached.
func setQueryParams(queryParams []QueryParam, url string) string {
	if len(queryParams) == 0 {
		return url
	}

	var builder strings.Builder
	builder.WriteString(url)

	for i, param := range queryParams {
		value := netUrl.QueryEscape(fmt.Sprintf("%v", param.Value))
		if i == 0 {
			builder.WriteString("?")
		} else {
			builder.WriteString("&")
		}
		builder.WriteString(param.Key)
		builder.WriteString("=")
		builder.WriteString(value)
	}

	return builder.String()
}
