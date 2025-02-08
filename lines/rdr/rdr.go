package lrdr

import (
	"bufio"
	"io"
	"iter"
	"os"
)

func ReaderToLines(rdr io.Reader) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		var s *bufio.Scanner = bufio.NewScanner(rdr)
		for s.Scan() {
			var line []byte = s.Bytes()
			if !yield(line) {
				return
			}
		}
	}
}

func StdinToLines() iter.Seq[[]byte] {
	return ReaderToLines(os.Stdin)
}
