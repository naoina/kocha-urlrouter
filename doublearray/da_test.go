package doublearray

import (
	"reflect"
	"testing"

	"github.com/naoina/kocha-urlrouter/testutil"
)

func Test_New(t *testing.T) {
	da := New()

	actual := reflect.TypeOf(da)
	expected := reflect.TypeOf(&DoubleArray{})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func Test_DoubleArray_Lookup(t *testing.T) {
	testutil.Test_URLRouter_Lookup(t, &DoubleArrayRouter{})
}

func Test_DoubleArray_Lookup_with_many_routes(t *testing.T) {
	testutil.Test_URLRouter_Lookup_with_many_routes(t, &DoubleArrayRouter{})
}

func Test_DoubleArray_Build(t *testing.T) {
	testutil.Test_URLRouter_Build(t, &DoubleArrayRouter{})
}
