// A URL router implemented by Ternary Search Tree.
package tst

import (
	"fmt"

	"github.com/naoina/kocha-urlrouter"
)

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
func (tst *TST) Build(records []urlrouter.Record) error {
	for _, record := range records {
		if err := tst.root.Add(record.Key, record.Value); err != nil {
			return err
		}
	}
	return nil
}

// node represents a node of TST.
type node struct {
	c            byte
	data         interface{}
	left         *node
	mid          *node
	right        *node
	paramNode    *node
	wildcardNode *node
	paramNames   []string
	isLeaf       bool
}

type nodeIndex struct {
	nd  *node
	idx int
}

func (nd *node) Find(path string, params []string) (*node, []string) {
	var nodes []nodeIndex
	for i := 0; i < len(path); i++ {
		if nd = nd.mid.find(path[i]); nd == nil {
			goto PARAMED_ROUTE
		}
		if nd.paramNode != nil || nd.wildcardNode != nil {
			nodes = append(nodes, nodeIndex{nd, i + 1})
		}
	}
	return nd, params
PARAMED_ROUTE:
	for i := len(nodes) - 1; i >= 0; i-- {
		nd, idx := nodes[i].nd, nodes[i].idx
		if nd.paramNode != nil {
			i := urlrouter.NextSeparator(path, idx)
			remaining, params := path[i:], append(params, path[idx:i])
			if nd, params := nd.paramNode.Find(remaining, params); nd != nil {
				return nd, params
			}
		}
		if nd.wildcardNode != nil {
			return nd.wildcardNode, append(params, path[idx:])
		}
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

func (nd *node) Add(path string, data interface{}) error {
	var paramNames []string
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
			paramNames = append(paramNames, remaining[:len(remaining)])
			nd.wildcardNode = &node{}
			nd = nd.wildcardNode
			i = len(path)
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
	if len(paramNames) > 0 {
		dups := make(map[string]bool)
		for _, name := range paramNames {
			if dups[name] {
				return fmt.Errorf("path parameter `%v` is duplicated in the key '%v'", name, path)
			}
			dups[name] = true
		}
	}
	nd.data, nd.paramNames, nd.isLeaf = data, paramNames, true
	return nil
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
