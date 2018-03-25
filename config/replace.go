package config

type Replace struct {
	Search  string
	Replace string
	Regexp  bool
}

func (r *Replace) IsValid() bool {
	if r.Search == "" {
		return false
	}

	if r.Replace == "" {
		return false
	}

	return true
}
