package livereloadproxy

import (
	"regexp"
)

func NewSourceReplace(search string, replace string) *SourceReplace {
	return &SourceReplace{
		searchRegexp: regexp.MustCompile(search),
		replace:      replace,
	}
}

type SourceReplace struct {
	searchRegexp *regexp.Regexp
	replace      string
}

func (s *SourceReplace) Replace(value string) string {
	return s.searchRegexp.ReplaceAllString(value, s.replace)
}
