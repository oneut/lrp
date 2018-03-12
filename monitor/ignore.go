package monitor

import (
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

func NewIgnore(ignoreString string) *Ignore {
	ignore := &Ignore{
		ignore: ignoreString,
	}
	ignore.Init()
	return ignore
}

type Ignore struct {
	ignore        string
	absIgnore     string
	ignoreGlob    glob.Glob
	absIgnoreGlob glob.Glob
}

func (i *Ignore) Init() {
	i.absIgnore = i.Abs()
	i.ignoreGlob = glob.MustCompile(i.GetIgnoreGlobString())
	i.absIgnoreGlob = glob.MustCompile(i.absIgnore)
}

func (i *Ignore) Match(path string) bool {
	// partial path check
	if i.Contains(path) {
		return true
	}

	if i.MatchGlob(path) {
		return true
	}

	// abstract path check
	if i.HasPrefixAbs(path) {
		return true
	}

	if i.MatchAbsGlob(path) {
		return true
	}

	return false
}

func (i *Ignore) GetIgnoreGlobString() string {
	return "*" + i.ignore
}

func (i *Ignore) Contains(path string) bool {
	return strings.Contains(path, i.ignore)
}

func (i *Ignore) HasPrefixAbs(path string) bool {
	return strings.HasPrefix(path, i.absIgnore)
}

func (i *Ignore) MatchGlob(path string) bool {
	return i.ignoreGlob.Match(path)
}

func (i *Ignore) MatchAbsGlob(path string) bool {
	return i.absIgnoreGlob.Match(path)
}

func (i *Ignore) Abs() string {
	if filepath.IsAbs(i.ignore) {
		return i.ignore
	}

	absPath, err := filepath.Abs(i.ignore)
	if err != nil {
		panic(err)
	}

	return absPath
}
