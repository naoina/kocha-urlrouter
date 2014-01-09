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

func Benchmark_Regexp_Build_100(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &RegexpRouter{}, 100)
}

func Benchmark_Regexp_Build_300(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &RegexpRouter{}, 300)
}

func Benchmark_Regexp_Build_700(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &RegexpRouter{}, 700)
}
