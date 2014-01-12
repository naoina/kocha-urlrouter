// A URL router implemented by Ternary Search Tree.
package tst

import "github.com/naoina/kocha-urlrouter"

// TST represents a URLRouter by Ternary Search Tree.
type TST struct {
	root *node
}

// New returns a new TST.
func New() *TST {
	return &TST{root: &node{}}
}

// Lookup returns result data of lookup from TST routing table by given path.
func (tst *TST) Lookup(path string) (data interface{}, params []urlrouter.Param) {
	nd, values := tst.root.Find(path, []string{})
	if nd == nil || !nd.isLeaf {
		return nil, nil
	}
	if len(values) > 0 {
		params = make([]urlrouter.Param, len(values))
		for i, v := range values {
			params[i] = urlrouter.Param{Name: nd.paramNames[i], Value: v}
		}
	}
	return nd.data, params
}

// Build builds TST routing table from records.
func (tst *TST) Build(records []*urlrouter.Record) error {
	for _, record := range records {
		tst.root.Add(record.Key, record.Value)
	}
	return nil
}

// node represents a node of TST.
type node struct {
	c              byte
	data           interface{}
	left           *node
	mid            *node
	right          *node
	paramNode      *node
	isWildcardNode bool
	paramNames     []string
	isLeaf         bool
}

func (nd *node) Find(path string, params []string) (*node, []string) {
	if path == "" {
		return nd, params
	}
	c, remaining := path[0], path[1:]
	if nd := nd.mid.find(c); nd != nil {
		if nd, params := nd.Find(remaining, params); nd != nil {
			return nd, params
		}
	}
	if nd.paramNode != nil {
		i := urlrouter.NextSeparator(path, 0)
		remaining, params = path[i:], append(params, path[:i])
		if nd, params = nd.paramNode.Find(remaining, params); nd != nil {
			return nd, params
		}
	}
	if nd.isWildcardNode {
		return nd, append(params, path)
	}
	return nil, nil
}

func (nd *node) find(c byte) *node {
	for nd != nil {
		switch {
		case nd.c > c:
			nd = nd.left
		case nd.c < c:
			nd = nd.right
		default: // n.c == c
			return nd
		}
	}
	return nil
}

func (nd *node) Add(path string, data interface{}) {
	var paramNames []string
LOOP:
	for i := 0; i < len(path); i++ {
		switch c, remaining := path[i], path[i+1:]; c {
		case urlrouter.ParamCharacter:
			next := urlrouter.NextSeparator(remaining, 0)
			paramNames = append(paramNames, remaining[:next])
			if nd.paramNode == nil {
				nd.paramNode = &node{}
			}
			i += next
			nd = nd.paramNode
		case urlrouter.WildcardCharacter:
			i := urlrouter.NextSeparator(remaining, 0)
			paramNames = append(paramNames, remaining[:i])
			nd.isWildcardNode = true
			break LOOP
		default:
			n := nd.mid.find(c)
			if n == nil {
				n = &node{c: c}
				if nd.mid == nil {
					nd.mid = n
				} else {
					nd.mid.add(n)
				}
			}
			nd = n
		}
	}
	nd.data, nd.paramNames, nd.isLeaf = data, paramNames, true
}

// add adds a node to leaf.
func (nd *node) add(n *node) {
	last := &nd
	for *last != nil {
		switch nd := *last; {
		case nd.c > n.c:
			last = &nd.left
		case nd.c < n.c:
			last = &nd.right
		default: // nd.c == n.c
			last = &nd.mid
		}
	}
	*last = n
}

// TSTRouter represents the Router of TST.
type TSTRouter struct{}

// New returns a new URLRouter that implemented by TST.
func (router *TSTRouter) New() urlrouter.URLRouter {
	return New()
}

func init() {
	urlrouter.Register("tst", &TSTRouter{})
}
