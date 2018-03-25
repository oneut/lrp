package livereloadproxy

import (
	"regexp"
)

func NewSourceRegexpReplacer(search string, replace string) *SourceRegexpReplacer {
	return &SourceRegexpReplacer{
		searchRegexp: regexp.MustCompile(search),
		replace:      replace,
	}
}

type SourceRegexpReplacer struct {
	searchRegexp *regexp.Regexp
	replace      string
}

func (s *SourceRegexpReplacer) Replace(value string) string {
	return s.searchRegexp.ReplaceAllString(value, s.replace)
}
