package tst

import (
	"reflect"
	"testing"

	"github.com/naoina/kocha-urlrouter/testutil"
)

func Test_New(t *testing.T) {
	re := New()

	actual := reflect.TypeOf(re)
	expected := reflect.TypeOf(&TST{})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func Test_TST_Lookup(t *testing.T) {
	testutil.Test_URLRouter_Lookup(t, New())
}

func Test_TST_Lookup_with_many_routes(t *testing.T) {
	testutil.Test_URLRouter_Lookup_with_many_routes(t, New())
}