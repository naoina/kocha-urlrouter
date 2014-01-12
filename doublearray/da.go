// A URL router implemented by Double-Array Trie.
package doublearray

import (
	"fmt"
	"sort"

	"github.com/naoina/kocha-urlrouter"
)

const (
	// Block size of array of BASE/CHECK of Double-Array.
	blockSize = 256
)

// baseCheck represents a BASE/CHECK node.
type baseCheck struct {
	base  int
	check int
}

// DoubleArray represents a URLRouter by Double-Array.
type DoubleArray struct {
	bc   []baseCheck
	node map[int]*node
}

// NewDoubleArray returns a new DoubleArray.
func New() *DoubleArray {
	da := &DoubleArray{
		bc:   newBaseCheckArray(blockSize),
		node: make(map[int]*node),
	}
	return da
}

// newBaseCheckArray returns a new slice of baseCheck with given size.
func newBaseCheckArray(size int) []baseCheck {
	bc := make([]baseCheck, size)
	for i := 0; i < len(bc); i++ {
		bc[i].check = -1
	}
	return bc
}

// Lookup returns result data of lookup from Double-Array routing table by given path.
func (da *DoubleArray) Lookup(path string) (data interface{}, params []urlrouter.Param) {
	nodes, idx, values := da.lookup(path, nil, 0)
	if nodes == nil {
		return nil, nil
	}
	nd := nodes[idx]
	if nd == nil {
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

// Build builds Double-Array routing table from records.
func (da *DoubleArray) Build(records []*urlrouter.Record) error {
	keys := sortedKeys(records)
	if err := da.build(keys, 0, 0); err != nil {
		return err
	}
	for _, record := range records {
		nodes, idx, names := da.lookup(record.Key, nil, 0)
		if nodes == nil {
			return fmt.Errorf("BUG: routing table could not be built correctly")
		}
		nodes[idx] = &node{}
		if len(names) > 0 {
			for i, name := range names {
				names[i] = name[1:] // truncate the meta character.
			}
			nodes[idx].paramNames = names
		}
		nodes[idx].data = record.Value
	}
	return nil
}

func (da *DoubleArray) lookup(path string, params []string, idx int) (map[int]*node, int, []string) {
	if path == "" {
		return da.node, idx - 1, params
	}
	c, remaining := path[0], path[1:]
	if next := nextIndex(da.bc[idx].base, c); da.bc[next].check == idx {
		if nodes, idx, params := da.lookup(remaining, params, next); nodes != nil {
			return nodes, idx, params
		}
	}
	if nd := da.node[idx]; nd != nil && nd.paramTree != nil {
		i := urlrouter.NextSeparator(path, 0)
		remaining, params = path[i:], append(params, path[:i])
		if nodes, idx, params := nd.paramTree.lookup(remaining, params, 0); nodes != nil {
			return nodes, idx, params
		}
	}
	if nd := da.node[idx]; nd != nil && nd.isWildcard {
		return da.node, idx - 1, append(params, path)
	}
	return nil, -1, nil
}

func (da *DoubleArray) build(routePaths []string, idx, depth int) error {
	base, siblings, err := da.arrange(routePaths, idx, depth)
	if err != nil {
		return err
	}
	for _, sib := range siblings {
		if !urlrouter.IsMetaChar(sib.c) {
			da.setCheck(nextIndex(base, sib.c), idx)
		}
	}
	for _, sib := range siblings {
		switch sib.c {
		case urlrouter.ParamCharacter:
			paths := routePaths[sib.start:sib.end]
			for i, path := range paths {
				paths[i] = path[urlrouter.NextSeparator(path, depth):]
			}
			da.node[idx] = &node{paramTree: New()}
			if err := da.node[idx].paramTree.build(paths, 0, 0); err != nil {
				return err
			}
		case urlrouter.WildcardCharacter:
			da.node[idx] = &node{isWildcard: true}
		default:
			if err := da.build(routePaths[sib.start:sib.end], nextIndex(base, sib.c), depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// setBase sets BASE.
func (da *DoubleArray) setBase(i, base int) {
	da.bc[i].base = base
}

// setCheck sets CHECK.
func (da *DoubleArray) setCheck(i, check int) {
	da.bc[i].check = check
}

// extendBaseCheckArray extends array of BASE/CHECK.
func (da *DoubleArray) extendBaseCheckArray() {
	da.bc = append(da.bc, newBaseCheckArray(blockSize)...)
}

// findEmptyIndex returns an index of unused BASE/CHECK node.
func (da *DoubleArray) findEmptyIndex(start int) int {
	i := start
	for ; i < len(da.bc); i++ {
		if da.bc[i].base == 0 && da.bc[i].check == -1 {
			break
		}
	}
	return i
}

// findBase returns good BASE.
func (da *DoubleArray) findBase(siblings []sibling, start int) (base int) {
	idx := start + 1
	firstChar := siblings[0].c
	for ; idx < len(da.bc); idx = da.findEmptyIndex(idx + 1) {
		base = nextIndex(idx, firstChar)
		i := 0
		for ; i < len(siblings); i++ {
			if urlrouter.IsMetaChar(siblings[i].c) {
				continue
			}
			if next := nextIndex(base, siblings[i].c); da.bc[next].base != 0 || da.bc[next].check != -1 {
				break
			}
		}
		if i == len(siblings) {
			return base
		}
	}
	da.extendBaseCheckArray()
	return nextIndex(idx, firstChar)
}

func (da *DoubleArray) arrange(keys []string, idx, depth int) (base int, siblings []sibling, err error) {
	siblings, err = makeSiblings(keys, depth)
	if err != nil {
		return -1, nil, err
	}
	if len(siblings) < 1 {
		return -1, nil, nil
	}
	base = da.findBase(siblings, idx)
	da.setBase(idx, base)
	return base, siblings, err
}

// node represents a node of Double-Array.
type node struct {
	data interface{}

	// Tree of path parameter.
	paramTree *DoubleArray

	// Names of path parameters.
	paramNames []string

	// Whether the wildcard node.
	isWildcard bool
}

// sibling represents an intermediate data of build for Double-Array.
type sibling struct {
	// An index of start of duplicated characters.
	start int

	// An index of end of duplicated characters.
	end int

	// A character of sibling.
	c byte
}

// nextIndex returns a next index of array of BASE/CHECK.
func nextIndex(base int, c byte) int {
	return base ^ int(c)
}

// makeSiblings returns slice of sibling from string keys.
func makeSiblings(keys []string, depth int) (sib []sibling, err error) {
	var (
		pc byte
		n  int
	)
	for i, key := range keys {
		if len(key) <= depth {
			continue
		}
		c := key[depth]
		switch {
		case pc < c:
			sib = append(sib, sibling{start: i, c: c})
		case pc == c:
			continue
		default:
			return nil, fmt.Errorf("BUG: routing table hasn't been sorted")
		}
		if n > 0 {
			sib[n-1].end = i
		}
		pc = c
		n++
	}
	if n == 0 {
		return nil, nil
	}
	sib[n-1].end = len(keys)
	return sib, nil
}

// sortedKeys returns sorted keys of records.
func sortedKeys(records []*urlrouter.Record) (keys []string) {
	keys = make([]string, len(records))
	for i, record := range records {
		keys[i] = record.Key
	}
	sort.Strings(keys)
	return keys
}

// DoubleArrayRouter represents the Router of Double-Array.
type DoubleArrayRouter struct{}

// New returns a new URLRouter that implemented by Double-Array.
func (router *DoubleArrayRouter) New() urlrouter.URLRouter {
	return New()
}

func init() {
	urlrouter.Register("doublearray", &DoubleArrayRouter{})
}
