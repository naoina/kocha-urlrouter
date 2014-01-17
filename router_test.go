package urlrouter

import (
	"reflect"
	"testing"
)

type testRouter struct {
	ur URLRouter
}

func (r *testRouter) New() URLRouter {
	return r.ur
}

type testURLRouter struct {
	name string
}

func (r *testURLRouter) Lookup(path string) (data interface{}, params []Param) {
	return nil, nil
}

func (r *testURLRouter) Build(records []Record) error {
	return nil
}

func Test_Register(t *testing.T) {
	defer func() {
		routers = make(map[string]Router)
	}()

	var actual interface{} = len(routers)
	var expected interface{} = 0
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
	router := &testRouter{}
	Register("testrouter", router)
	actual = routers
	expected = map[string]Router{"testrouter": router}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func Test_NewURLRouter(t *testing.T) {
	defer func() {
		routers = make(map[string]Router)
	}()
	router1, router2 := &testURLRouter{name: "1"}, &testURLRouter{name: "2"}
	routers["router1"] = &testRouter{router1}
	routers["router2"] = &testRouter{router2}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expect error, but nil")
			}
		}()
		NewURLRouter("")
	}()

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expect error, but nil")
			}
		}()
		NewURLRouter("missing")
	}()

	actual := NewURLRouter("router1")
	expected := router1
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	actual = NewURLRouter("router2")
	expected = router2
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func Test_NewRecord(t *testing.T) {
	actual := NewRecord("testkey", 100)
	expected := Record{Key: "testkey", Value: 100}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}
