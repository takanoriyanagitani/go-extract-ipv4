package main

import (
	"context"
	"iter"
	"log"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
	ff "github.com/takanoriyanagitani/go-extract-ipv4/filter/fine"
	lr "github.com/takanoriyanagitani/go-extract-ipv4/lines/rdr"
	lw "github.com/takanoriyanagitani/go-extract-ipv4/lines/wtr"
	. "github.com/takanoriyanagitani/go-extract-ipv4/util"
)

var cfilt ei.CoarseFilter = ei.CoarseFilterDefault

var extpat ff.ExtPattern = ff.ExtPatternDefault

var extip4 ff.ExtractIps4 = extpat.ToExtractIps4()

var ffilt ff.FineFilter = extip4.ToFineFilter(cfilt)

var stdin2lines IO[iter.Seq[[]byte]] = Of(lr.StdinToLines())

var filtered IO[iter.Seq2[[]byte, error]] = Bind(
	stdin2lines,
	ffilt.ToFoundLines,
)

var lines2stdout func(iter.Seq2[[]byte, error]) IO[Void] = lw.LineWriterStdout

var stdin2lines2filtered2stdout IO[Void] = Bind(
	filtered,
	lines2stdout,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return stdin2lines2filtered2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
