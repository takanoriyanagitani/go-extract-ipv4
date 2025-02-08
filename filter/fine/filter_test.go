package fine_test

import (
	"context"
	"testing"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
	ff "github.com/takanoriyanagitani/go-extract-ipv4/filter/fine"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	t.Run("ExtPatternDefault", func(t *testing.T) {
		t.Parallel()

		var pat ff.ExtPattern = ff.ExtPatternDefault

		var cfilt ei.CoarseFilter = ei.CoarseFilterDefault

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var ex4 ff.ExtractIps4 = pat.ToExtractIps4()
			var ffilt ff.FineFilter = ex4.ToFineFilter(cfilt)

			var empty []byte

			res, e := ffilt(empty)(context.Background())
			if nil != e {
				t.Fatalf("unexpected error: %v\n", e)
			}

			if res {
				t.Fatalf("must not match\n")
			}
		})

		t.Run("localhost", func(t *testing.T) {
			t.Parallel()

			var ex4 ff.ExtractIps4 = pat.ToExtractIps4()
			var ffilt ff.FineFilter = ex4.ToFineFilter(cfilt)

			var sample []byte = []byte("helo 127.0.0.1 world")

			res, e := ffilt(sample)(context.Background())
			if nil != e {
				t.Fatalf("unexpected error: %v\n", e)
			}

			if !res {
				t.Fatalf("must match\n")
			}
		})
	})
}
