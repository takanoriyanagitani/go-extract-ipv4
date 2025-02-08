package wtr

import (
	"bufio"
	"context"
	"io"
	"iter"
	"os"

	. "github.com/takanoriyanagitani/go-extract-ipv4/util"
)

func LinesToWriter(
	ctx context.Context,
	lines iter.Seq2[[]byte, error],
	wtr io.Writer,
) error {
	for line, e := range lines {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if nil != e {
			return e
		}

		_, e := wtr.Write(line)
		if nil != e {
			return e
		}

		_, e = wtr.Write([]byte("\n"))
		if nil != e {
			return e
		}
	}
	return nil
}

func WriterToLineWriter(wtr io.Writer) func(iter.Seq2[[]byte, error]) IO[Void] {
	return func(lines iter.Seq2[[]byte, error]) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			return Empty, LinesToWriter(ctx, lines, wtr)
		}
	}
}

func LinesToStdout(lines iter.Seq2[[]byte, error]) IO[Void] {
	var bw *bufio.Writer = bufio.NewWriter(os.Stdout)
	return Bind(
		WriterToLineWriter(bw)(lines),
		Lift(func(_ Void) (Void, error) {
			return Empty, bw.Flush()
		}),
	)
}

var LineWriterStdout func(iter.Seq2[[]byte, error]) IO[Void] = LinesToStdout
