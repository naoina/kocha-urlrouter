// A URL router implemented by Regular-Expression.
package regexp

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/naoina/kocha-urlrouter"
)

const (
	defaultParamRegexpStr = `[\w-]+`
)

var (
	pathRegexp     = regexp.MustCompile(`/(?:(:[\w-]+)|(\*[\w-/]+)|[\w-]*)`)
	paramRegexpStr = map[byte]string{
		':': `[\w-]+`,
		'*': `[\w-/.]+`,
	}
)

// Regexp represents a URLRouter by Regular-Expression.
type Regexp struct {
	routes []*route
}

// New returns a new Regexp.
func New() *Regexp {
	return &Regexp{}
}

// Lookup returns result data of lookup from regexp routing table by given path.
func (re *Regexp) Lookup(path string) (data interface{}, params []urlrouter.Param) {
	for _, nd := range re.routes {
		matchesBase := nd.regexp.FindStringSubmatch(path)
		if len(matchesBase) < 1 {
			continue
		}
		subexpNames := nd.regexp.SubexpNames()[1:]
		if matches := matchesBase[1:]; len(matches) > 0 {
			params = make([]urlrouter.Param, len(matches))
			for i := 0; i < len(matches); i++ {
				params[i] = urlrouter.Param{Name: subexpNames[i], Value: matches[i]}
			}
		}
		return nd.data, params
	}
	return nil, nil
}

// Build builds regexp routing table from records.
func (re *Regexp) Build(records []urlrouter.Record) error {
	re.routes = make([]*route, len(records))
	for i, record := range records {
		route, err := build(record.Key, record.Value)
		if err != nil {
			return err
		}
		re.routes[i] = route
	}
	return nil
}

func build(path string, data interface{}) (*route, error) {
	var buf bytes.Buffer
	dups := make(map[string]bool)
	for _, paths := range pathRegexp.FindAllStringSubmatch(path, -1) {
		name := paths[1] + paths[2]
		if name == "" {
			// don't have path parameters.
			buf.WriteString(regexp.QuoteMeta(paths[0]))
			continue
		}
		var pathReStr string
		if pathReStr = paramRegexpStr[name[0]]; pathReStr == "" {
			pathReStr = defaultParamRegexpStr
		} else {
			if dups[name] {
				return nil, fmt.Errorf("path parameter `%v` is duplicated in the key '%v'", name, path)
			}
			dups[name] = true
			name = name[1:] // truncate a meta character.
		}
		buf.WriteString(fmt.Sprintf(`/(?P<%s>%s)`, regexp.QuoteMeta(name), pathReStr))
	}
	reg, err := regexp.Compile(fmt.Sprintf(`^%s$`, buf.String()))
	if err != nil {
		return nil, err
	}
	return &route{regexp: reg, data: data}, nil
}

// route represents a regexp route.
type route struct {
	regexp *regexp.Regexp
	data   interface{}
}

// RegexpRouter represents the Router of Regular-Expression.
type RegexpRouter struct{}

// New returns a new URLRouter that implemented by Regular-Expression.
func (router *RegexpRouter) New() urlrouter.URLRouter {
	return New()
}

func init() {
	urlrouter.Register("regexp", &RegexpRouter{})
}
