package exipv4

import (
	"regexp"
)

type CoarseFilterResult bool

const (
	CoarseFilterNoMatch CoarseFilterResult = false
	CoarseFilterMayIpIV CoarseFilterResult = true
)

type CoarseFilter func([]byte) CoarseFilterResult

type Pattern struct {
	*regexp.Regexp
}

func (p Pattern) ToCoarseFilter() CoarseFilter {
	return func(line []byte) CoarseFilterResult {
		var b bool = p.Regexp.Match(line)
		return CoarseFilterResult(b)
	}
}

var PatternDefault Pattern = Pattern{
	Regexp: regexp.MustCompile(`\d\.`),
}

var CoarseFilterDefault CoarseFilter = PatternDefault.ToCoarseFilter()

type Ip4 [4]byte
