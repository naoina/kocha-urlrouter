package testutil

import (
	"fmt"
	"testing"

	"github.com/naoina/kocha-urlrouter"
)

func Benchmark_URLRouter_Lookup(b *testing.B, router urlrouter.URLRouter, n int) {
	b.StopTimer()
	records := makeTestRecords(n)
	if err := router.Build(records); err != nil {
		b.Fatal(err)
	}
	record := pickTestRecord(records)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if r, _ := router.Lookup(record.Key); r != record.Value {
			b.Fail()
		}
	}
}

func Benchmark_URLRouter_Build(b *testing.B, router urlrouter.Router, n int) {
	b.StopTimer()
	records := makeTestRecords(n)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r := router.New()
		if err := r.Build(records); err != nil {
			b.Fatal(err)
		}
	}
}

func makeTestRecords(n int) []urlrouter.Record {
	records := make([]urlrouter.Record, n)
	for i := 0; i < n; i++ {
		records[i] = urlrouter.NewRecord("/"+RandomString(50), fmt.Sprintf("testroute%d", i))
	}
	return records
}

func pickTestRecord(records []urlrouter.Record) urlrouter.Record {
	return records[len(records)/2]
}
