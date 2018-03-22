package config

var defaultSourceScheme string = "http"

type Source struct {
	Scheme   string
	Host     string
	Replaces []Replace
}

func (s *Source) GetScheme() string {
	if s.Scheme == "" {
		return defaultSourceScheme
	}

	return s.Scheme
}
