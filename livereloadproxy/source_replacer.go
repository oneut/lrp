package livereloadproxy

import (
	"github.com/oneut/lrp/config"
)

func NewSourceReplacer(replaceConfig config.Replace) SourceReplacer {
	if replaceConfig.Regexp {
		return NewSourceRegexpReplacer(replaceConfig.Search, replaceConfig.Replace)
	}

	return NewSourceStringReplacer(replaceConfig.Search, replaceConfig.Replace)
}

type SourceReplacer interface {
	Replace(string) string
}
