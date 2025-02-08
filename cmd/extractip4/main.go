package main

import (
	"context"
	"iter"
	"log"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
	ff "github.com/takanoriyanagitani/go-extract-ipv4/filter/fine"
	iw "github.com/takanoriyanagitani/go-extract-ipv4/ipv4/wtr"
	lr "github.com/takanoriyanagitani/go-extract-ipv4/lines/rdr"
	. "github.com/takanoriyanagitani/go-extract-ipv4/util"
)

var cfilt ei.CoarseFilter = ei.CoarseFilterDefault

var extpat ff.ExtPattern = ff.ExtPatternDefault

var extip4 ff.ExtractIps4 = extpat.ToExtractIps4()

var stdin2lines IO[iter.Seq[[]byte]] = Of(lr.StdinToLines())

var ips IO[iter.Seq2[[]ei.Ip4, error]] = Bind(
	stdin2lines,
	func(lines iter.Seq[[]byte]) IO[iter.Seq2[[]ei.Ip4, error]] {
		return extip4.ToFoundIps(lines, cfilt)
	},
)

var ips2stdoutRaw func(iter.Seq2[[]ei.Ip4, error]) IO[Void] = iw.
	StdoutToIpsWriterRaw

var stdin2lines2ips2stdoutRaw IO[Void] = Bind(
	ips,
	ips2stdoutRaw,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2lines2ips2stdoutRaw(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
