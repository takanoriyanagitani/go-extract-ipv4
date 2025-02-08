package wtr

import (
	"bufio"
	"context"
	"io"
	"iter"
	"os"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
	. "github.com/takanoriyanagitani/go-extract-ipv4/util"
)

func WriteIpRaw(ip ei.Ip4, wtr io.Writer) error {
	var s []byte = ip[:]
	_, e := wtr.Write(s)
	return e
}

func WriterToIpsWriterRaw(
	wtr io.Writer,
) func(iter.Seq2[[]ei.Ip4, error]) IO[Void] {
	return func(ips iter.Seq2[[]ei.Ip4, error]) IO[Void] {
		return func(ctx context.Context) (Void, error) {
			for items, e := range ips {
				select {
				case <-ctx.Done():
					return Empty, ctx.Err()
				default:
				}

				if nil != e {
					return Empty, e
				}

				for _, ipv4 := range items {
					e := WriteIpRaw(ipv4, wtr)
					if nil != e {
						return Empty, e
					}
				}
			}

			return Empty, nil
		}
	}
}

func StdoutToIpsWriterRaw(lines iter.Seq2[[]ei.Ip4, error]) IO[Void] {
	var bw *bufio.Writer = bufio.NewWriter(os.Stdout)
	return Bind(
		WriterToIpsWriterRaw(bw)(lines),
		Lift(func(_ Void) (Void, error) {
			return Empty, bw.Flush()
		}),
	)
}
