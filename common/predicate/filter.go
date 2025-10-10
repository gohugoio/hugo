package predicate

import (
	"iter"
	"strings"

	hglob "github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gobwas/glob"
)

type Filter[T any] struct {
	exclude P[T]
	include P[T]
}

// TODO1 remove this filter.
func (f Filter[T]) ShouldExclude(v T) bool {
	return f.ShouldExcludeFine(v)
}

func (f Filter[T]) ShouldExcludeFine(v T) bool {
	if f.exclude != nil && f.exclude(v) {
		return true
	}
	return f.include != nil && !f.include(v)
}

func NewFilter[T any](include, exclude P[T]) Filter[T] {
	return Filter[T]{exclude: exclude, include: include}
}

func NewStringFilterFromGlobs(patterns []string, getGlob func(pattern string) (glob.Glob, error)) (Filter[string], error) {
	var filter Filter[string]
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, hglob.NegationPrefix) {
			p = p[2:]
			g, err := getGlob(p)
			if err != nil {
				return filter, err
			}
			fn := func(s string) bool {
				return g.Match(s)
			}
			filter.exclude = filter.exclude.Or(fn)
		} else {
			g, err := getGlob(p)
			if err != nil {
				return filter, err
			}
			fn := func(s string) bool {
				return g.Match(s)
			}
			filter.include = filter.include.Or(fn)
		}
	}

	return filter, nil
}

type IndexMatcher interface {
	IndexMatch(filter Filter[string]) (iter.Seq[int], error)
}
