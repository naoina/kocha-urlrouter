package urlrouter

import "fmt"

const (
	ParamCharacter    = ':'
	WildcardCharacter = '*'
)

var routers map[string]Router

// URLRouter is an interface that must be implemented by a URL router.
type URLRouter interface {
	// Lookup returns data and path parameters that associated with path.
	// params is a slice of the Param that arranged in the order in which parameters appeared.
	// e.g. when built routing path is "/path/:id/:name" and given path is "/path/to/1/alice". params order is [{"id": "1"}, {"name": "alice"}], not [{"name": "alice"}, {"id": "1"}].
	// If failed to lookup, data will be nil.
	Lookup(path string) (data interface{}, params []Param)

	// Build builds URL router from records.
	Build(records []Record) error
}

// param represents a name and value of path parameter.
type Param struct {
	Name  string
	Value string
}

// Router is an interface of factory of URLRouter.
type Router interface {
	// New returns a new URLRouter.
	New() URLRouter
}

// Register registers a Router with name.
func Register(name string, router Router) {
	routers[name] = router
}

// NewURLRouter returns the URLRouter with the specified name.
func NewURLRouter(name string) URLRouter {
	router, exists := routers[name]
	if !exists {
		panic(fmt.Errorf("Router named `%v` is not registered", name))
	}
	return router.New()
}

// Record represents a record data for a router construction.
type Record struct {
	// Key for a router construction.
	Key string

	// Result value for Key.
	Value interface{}
}

// NewRecord returns a new Record.
func NewRecord(key string, value interface{}) Record {
	return Record{
		Key:   key,
		Value: value,
	}
}

func init() {
	routers = make(map[string]Router)
}
