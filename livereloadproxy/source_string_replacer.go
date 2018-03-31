package livereloadproxy

import "strings"

func NewSourceStringReplacer(search string, replace string) *SourceStringReplacer {
	return &SourceStringReplacer{
		search:  search,
		replace: replace,
	}
}

type SourceStringReplacer struct {
	search  string
	replace string
}

func (s *SourceStringReplacer) Replace(value string) string {
	return strings.Replace(value, s.search, s.replace, -1)
}
