package fine

import (
	"context"
	"iter"
	"net/netip"
	"regexp"
	"strings"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
	. "github.com/takanoriyanagitani/go-extract-ipv4/util"
)

type FilterResult bool

const (
	FilterResultNotFound FilterResult = false
	FilterResultFound    FilterResult = true
)

type FineFilter func([]byte) IO[FilterResult]

func (f FineFilter) ToFoundLines(
	lines iter.Seq[[]byte],
) IO[iter.Seq2[[]byte, error]] {
	return func(ctx context.Context) (iter.Seq2[[]byte, error], error) {
		return func(yield func([]byte, error) bool) {
			for line := range lines {
				select {
				case <-ctx.Done():
					yield(nil, ctx.Err())
					return
				default:
				}

				res, e := f(line)(ctx)
				if nil != e {
					yield(nil, e)
					return
				}

				if res {
					if !yield(line, nil) {
						return
					}
				}
			}
		}, nil
	}
}

type ExtractIps4 func([]byte) IO[[]ei.Ip4]

func (e ExtractIps4) FindIps(
	line []byte,
	res ei.CoarseFilterResult,
) IO[[]ei.Ip4] {
	return func(ctx context.Context) ([]ei.Ip4, error) {
		switch res {
		case ei.CoarseFilterMayIpIV:
			return e(line)(ctx)
		default:
			return nil, nil
		}
	}
}

func (e ExtractIps4) ToFineFilter(c ei.CoarseFilter) FineFilter {
	return func(line []byte) IO[FilterResult] {
		return Bind(
			e.FindIps(line, c(line)),
			Lift(func(ips []ei.Ip4) (FilterResult, error) {
				switch 0 < len(ips) {
				case true:
					return true, nil
				default:
					return false, nil
				}
			}),
		)
	}
}

func (e ExtractIps4) ToFoundIps(
	lines iter.Seq[[]byte],
	cfilt ei.CoarseFilter,
) IO[iter.Seq2[[]ei.Ip4, error]] {
	return func(ctx context.Context) (iter.Seq2[[]ei.Ip4, error], error) {
		return func(yield func([]ei.Ip4, error) bool) {
			for line := range lines {
				select {
				case <-ctx.Done():
					yield(nil, ctx.Err())
					return
				default:
				}

				var cres ei.CoarseFilterResult = cfilt(line)
				ips, err := e.FindIps(line, cres)(ctx)

				if !yield(ips, err) {
					return
				}
			}
		}, nil
	}
}

type ExtPattern struct {
	*regexp.Regexp
}

func (p ExtPattern) ToExtractIps4() ExtractIps4 {
	var buf []ei.Ip4
	var bld strings.Builder
	return func(line []byte) IO[[]ei.Ip4] {
		return func(_ context.Context) ([]ei.Ip4, error) {
			buf = buf[:0]

			var candidates [][]byte = p.Regexp.FindAll(line, -1)
			for _, ipc := range candidates {
				bld.Reset()
				_, _ = bld.Write(ipc) // error is always nil or OOM
				var s string = bld.String()

				addr, e := netip.ParseAddr(s)
				if nil != e {
					continue
				}

				if !addr.Is4() {
					continue
				}

				var ip []byte = addr.AsSlice()
				var ip4 ei.Ip4
				copy(ip4[:], ip)

				buf = append(buf, ip4)
			}

			return buf, nil
		}
	}
}

var ExtPatternDefault ExtPattern = ExtPattern{
	Regexp: regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
}
