package testutil

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/naoina/kocha-urlrouter"
)

func Test_URLRouter_Lookup(t *testing.T, router urlrouter.URLRouter) {
	testcases := []struct {
		path   string
		params map[string]string
		record *urlrouter.Record
	}{
		{"/", nil, &urlrouter.Record{"/", "testroute0"}},
		{"/path/to/route", nil, &urlrouter.Record{"/path/to/route", "testroute1"}},
		{"/path/to/other", nil, &urlrouter.Record{"/path/to/other", "testroute2"}},
		{"/path/to/route/a", nil, &urlrouter.Record{"/path/to/route/a", "testroute3"}},
		{"/path/to/hoge", map[string]string{"param": "hoge"}, &urlrouter.Record{"/path/to/:param", "testroute4"}},
		{"/path/to/wildcard/some/params", map[string]string{"routepath": "some/params"}, &urlrouter.Record{"/path/to/wildcard/*routepath", "testroute5"}},
		{"/path/to/o1/o2", map[string]string{"param1": "o1", "param2": "o2"}, &urlrouter.Record{"/path/to/:param1/:param2", "testroute6"}},
		{"/path/to/p1/sep/p2", map[string]string{"param1": "p1", "param2": "p2"}, &urlrouter.Record{"/path/to/:param1/sep/:param2", "testroute7"}},
		{"/2014/01/06", map[string]string{"year": "2014", "month": "01", "day": "06"}, &urlrouter.Record{"/:year/:month/:day", "testroute8"}},
		{"/user/777", map[string]string{"id": "777"}, &urlrouter.Record{"/user/:id", "testroute9"}},
		{"/a/to/b/p1/some/wildcard/params", map[string]string{"param": "p1", "routepath": "some/wildcard/params"}, &urlrouter.Record{"/a/to/b/:param/*routepath", "testroute10"}},
	}
	var records []*urlrouter.Record
	for _, testcase := range testcases {
		records = append(records, testcase.record)
	}
	if err := router.Build(records); err != nil {
		t.Fatal(err)
	}

	for _, testcase := range testcases {
		var actual, expected interface{}
		actual, params := router.Lookup(testcase.path)
		expected = testcase.record.Value
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

func Test_URLRouter_Lookup_with_many_routes(t *testing.T, router urlrouter.URLRouter) {
	n := 1000
	rand.Seed(time.Now().UnixNano())
	records := make([]*urlrouter.Record, n)
	for i := 0; i < n; i++ {
		records[i] = &urlrouter.Record{"/" + RandomString(rand.Intn(50)+10), fmt.Sprintf("route%d", i)}
	}
	if err := router.Build(records); err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		data, params := router.Lookup(record.Key)

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
