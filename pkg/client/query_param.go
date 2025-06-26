package client

import (
	"fmt"
	netUrl "net/url"

	"github.com/i9si-sistemas/stringx"
)

// QueryParam represents a query parameter in the URL of an HTTP request.
type QueryParam struct {
	Data // Embeds the Data struct to represent the key-value pair for the query parameter.
}

// SetQueryParams appends query parameters to the given URL.
// It returns the URL with the query parameters attached.
func SetQueryParams(queryParams []QueryParam, url string) string {
	if len(queryParams) == 0 {
		return url
	}

	builder := stringx.String(url).Builder()

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
