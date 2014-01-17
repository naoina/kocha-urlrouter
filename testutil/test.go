package testutil

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/naoina/kocha-urlrouter"
)

func routes() []urlrouter.Record {
	return []urlrouter.Record{
		{"/", "testroute0"},
		{"/path/to/route", "testroute1"},
		{"/path/to/other", "testroute2"},
		{"/path/to/route/a", "testroute3"},
		{"/path/to/:param", "testroute4"},
		{"/path/to/wildcard/*routepath", "testroute5"},
		{"/path/to/:param1/:param2", "testroute6"},
		{"/path/to/:param1/sep/:param2", "testroute7"},
		{"/:year/:month/:day", "testroute8"},
		{"/user/:id", "testroute9"},
		{"/a/to/b/:param/*routepath", "testroute10"},
	}
}

func Test_URLRouter_Lookup(t *testing.T, router urlrouter.Router) {
	type testcase struct {
		path   string
		value  interface{}
		params []urlrouter.Param
	}
	runTest := func(records []urlrouter.Record, testcases []*testcase) {
		r := router.New()
		if err := r.Build(records); err != nil {
			t.Fatal(err)
		}

		for _, testcase := range testcases {
			var actual, expected interface{}
			actual, params := r.Lookup(testcase.path)
			expected = testcase.value
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Expect %v, but %v", expected, actual)
			}

			actual = params
			expected = testcase.params
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("Expect %v, but %v", expected, actual)
			}
		}
	}
	testcases := []*testcase{
		{"/", "testroute0", nil},
		{"/path/to/route", "testroute1", nil},
		{"/path/to/other", "testroute2", nil},
		{"/path/to/route/a", "testroute3", nil},
		{"/path/to/hoge", "testroute4", []urlrouter.Param{{"param", "hoge"}}},
		{"/path/to/wildcard/some/params", "testroute5", []urlrouter.Param{{"routepath", "some/params"}}},
		{"/path/to/o1/o2", "testroute6", []urlrouter.Param{{"param1", "o1"}, {"param2", "o2"}}},
		{"/path/to/p1/sep/p2", "testroute7", []urlrouter.Param{{"param1", "p1"}, {"param2", "p2"}}},
		{"/2014/01/06", "testroute8", []urlrouter.Param{{"year", "2014"}, {"month", "01"}, {"day", "06"}}},
		{"/user/777", "testroute9", []urlrouter.Param{{"id", "777"}}},
		{"/a/to/b/p1/some/wildcard/params", "testroute10", []urlrouter.Param{{"param", "p1"}, {"routepath", "some/wildcard/params"}}},
		{"/missing", nil, nil},
	}
	runTest(routes(), testcases)

	records := []urlrouter.Record{
		{"/", "testroute0"},
		{"/:b", "testroute1"},
		{"/*wildcard", "testroute2"},
	}
	testcases = []*testcase{
		{"/", "testroute0", nil},
		{"/true", "testroute1", []urlrouter.Param{{"b", "true"}}},
		{"/foo/bar", "testroute2", []urlrouter.Param{{"wildcard", "foo/bar"}}},
	}
	runTest(records, testcases)
}

func Test_URLRouter_Lookup_with_many_routes(t *testing.T, router urlrouter.Router) {
	n := 1000
	rand.Seed(time.Now().UnixNano())
	records := make([]urlrouter.Record, n)
	for i := 0; i < n; i++ {
		records[i] = urlrouter.Record{"/" + RandomString(rand.Intn(50)+10), fmt.Sprintf("route%d", i)}
	}
	r := router.New()
	if err := r.Build(records); err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		data, params := r.Lookup(record.Key)

		var actual interface{} = data
		var expected interface{} = record.Value
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}

		actual = len(params)
		expected = 0
		if actual != expected {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}
}

func Test_URLRouter_Build(t *testing.T, router urlrouter.Router) {
	// test for duplicate name of path parameters.
	func() {
		r := router.New()
		if err := r.Build([]urlrouter.Record{
			{"/:user/:id/:id", "testroute0"},
			{"/:user/:user/:id", "testroute0"},
		}); err == nil {
			t.Errorf("no error returned by duplicate name of path parameters")
		}
	}()
}
