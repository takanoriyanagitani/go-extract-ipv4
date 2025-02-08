package exipv4_test

import (
	"testing"

	ei "github.com/takanoriyanagitani/go-extract-ipv4"
)

func TestExIp4(t *testing.T) {
	t.Parallel()

	t.Run("CoarseFilterDefault", func(t *testing.T) {
		t.Parallel()

		var filt ei.CoarseFilter = ei.CoarseFilterDefault

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var empty []byte
			var rslt ei.CoarseFilterResult = filt(empty)
			if ei.CoarseFilterNoMatch != rslt {
				t.Fatalf("must not match\n")
			}
		})

		t.Run("may contain ipv4", func(t *testing.T) {
			t.Parallel()

			var sample []byte = []byte("helo 127.0.0.1 world")
			var rslt ei.CoarseFilterResult = filt(sample)
			if ei.CoarseFilterNoMatch == rslt {
				t.Fatalf("must match\n")
			}
		})

		t.Run("no ipv4", func(t *testing.T) {
			t.Parallel()

			var sample []byte = []byte("helo 127 0 0 1 world")
			var rslt ei.CoarseFilterResult = filt(sample)
			if ei.CoarseFilterNoMatch != rslt {
				t.Fatalf("must not match\n")
			}
		})
	})
}
