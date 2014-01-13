package regexp

import (
	"reflect"
	"testing"

	"github.com/naoina/kocha-urlrouter/testutil"
)

func Test_New(t *testing.T) {
	re := New()

	actual := reflect.TypeOf(re)
	expected := reflect.TypeOf(&Regexp{})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func Test_Regexp_Lookup(t *testing.T) {
	testutil.Test_URLRouter_Lookup(t, &RegexpRouter{})
}

func Test_Regexp_Lookup_with_many_routes(t *testing.T) {
	testutil.Test_URLRouter_Lookup_with_many_routes(t, &RegexpRouter{})
}

func Test_Regexp_Build(t *testing.T) {
	testutil.Test_URLRouter_Build(t, &RegexpRouter{})
}
