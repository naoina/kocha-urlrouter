package urlrouter

import (
	"reflect"
	"testing"
)

func Test_NextSeparator(t *testing.T) {
	for _, testcase := range []struct {
		path     string
		start    int
		expected interface{}
	}{
		{"/path/to/route", 0, 0},
		{"/path/to/route", 1, 5},
		{"/path/to/route", 9, 14},
		{"/path.html", 1, 5},
		{"/foo/bar.html", 1, 4},
		{"/foo/bar.html/baz.png", 5, 8},
		{"/foo/bar.html/baz.png", 10, 13},
	} {
		actual := NextSeparator(testcase.path, testcase.start)
		expected := testcase.expected
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("path = %q, start = %v expect %v, but %v", testcase.path, testcase.start, expected, actual)
		}
	}
}

func Test_IsMetaChar(t *testing.T) {
	for _, c := range []byte{':', '*'} {
		if !IsMetaChar(c) {
			t.Errorf("Expect %q is meta charcter, but isn't", c)
		}
	}
	for c := byte(0); c < 0xff && c != ':' && c != '*'; c++ {
		if IsMetaChar(c) {
			t.Errorf("Expect %q is not meta character, but isn't", c)
		}
	}
}

func Test_ParamNames(t *testing.T) {
	for path, expected := range map[string][]string{
		"/:a":    {":a"},
		"/:b":    {":b"},
		"/:a/:b": {":a", ":b"},
		"/:ab":   {":ab"},
		"/*w":    {"*w"},
		"/*w/:p": {"*w", ":p"},
	} {
		actual := ParamNames(path)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%q expects %q, but %q", path, expected, actual)
		}
	}
}
