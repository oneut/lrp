package livereloadproxy

import (
	"net/http"
	"time"
)

type RetryTransport struct {
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		res, err := http.DefaultTransport.RoundTrip(req)
		if err == nil {
			return res, nil
		}

		// Retry
		time.Sleep(500 * time.Millisecond)
	}
}
