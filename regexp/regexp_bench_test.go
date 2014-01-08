package regexp

import (
	"testing"

	"github.com/naoina/kocha-urlrouter/testutil"
)

func Benchmark_Regexp_Lookup_100(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 100)
}

func Benchmark_Regexp_Lookup_300(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 300)
}

func BenchmarkRegexp_Lookup_700(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 700)
}
