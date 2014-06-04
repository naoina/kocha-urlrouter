package doublearray

import (
	"testing"

	"github.com/naoina/kocha-urlrouter/testutil"
)

func Benchmark_DoubleArray_Lookup_100(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 100)
}

func Benchmark_DoubleArray_Lookup_300(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 300)
}

func Benchmark_DoubleArray_Lookup_700(b *testing.B) {
	testutil.Benchmark_URLRouter_Lookup(b, New(), 700)
}

func Benchmark_DoubleArray_Build_100(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &DoubleArrayRouter{}, 100)
}

func Benchmark_DoubleArray_Build_300(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &DoubleArrayRouter{}, 300)
}

func Benchmark_DoubleArray_Build_700(b *testing.B) {
	testutil.Benchmark_URLRouter_Build(b, &DoubleArrayRouter{}, 700)
}
