package rmq

import (
	"net/url"
	"strings"
)

func GenerateUTMURL(baseURL, source, medium, campaign string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("utm_source", source)
	query.Set("utm_medium", medium)
	query.Set("utm_campaign", campaign)

	u.RawQuery = query.Encode()

	return u.String(), nil
}

func GetEventFromParams(queryParams string) EventType {
	if strings.TrimSpace(queryParams) == "" {
		return Read
	}

	queryParamsMap, err := url.ParseQuery(queryParams)
	if err != nil {
		return Read
	}

	utmMedium := queryParamsMap.Get("utm_medium")

	switch {
	case utmMedium == "email":
		return LinkClick
	default:
		return Read
	}
}
