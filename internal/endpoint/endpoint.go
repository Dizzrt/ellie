package endpoint

import "net/url"

func New(scheme, host string) *url.URL {
	return &url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

func Parse(endpoints []string, scheme string) (string, error) {
	for _, e := range endpoints {
		u, err := url.Parse(e)
		if err != nil {
			return "", err
		}

		if u.Scheme == scheme {
			return u.Host, nil
		}
	}

	return "", nil
}

func Scheme(scheme string, isSecure bool) string {
	if isSecure {
		return scheme + "s"
	}

	return scheme
}
