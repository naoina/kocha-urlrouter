// A URL router implemented by Double-Array Trie.
package doublearray

import (
	"fmt"
	"sort"

	"github.com/naoina/kocha-urlrouter"
)

const (
	// Block size of array of BASE/CHECK of Double-Array.
	blockSize         = 256
	paramCharacter    = ':'
	wildcardCharacter = '*'
)

// DoubleArray represents a URLRouter by Double-Array.
type DoubleArray struct {
	// BASE/CHECK array.
	bc []*node
}

// NewDoubleArray returns a new DoubleArray.
func New() *DoubleArray {
	da := &DoubleArray{bc: make([]*node, blockSize)}
	da.bc[0] = newNode()
	return da
}

// Lookup returns result data of lookup from Double-Array routing table by given path.
func (da *DoubleArray) Lookup(path string) (data interface{}, params map[string]string) {
	nd, values := da.lookup(path, []string{}, 0)
	if nd == nil || nd.data == nil {
		return nil, nil
	}
	if len(values) > 0 {
		params = make(map[string]string, len(values))
		for i, v := range values {
			params[nd.paramNames[i]] = v
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
		nd, names := da.lookup(record.Key, nil, 0)
		if nd == nil {
			return fmt.Errorf("BUG: routing table could not be built correctly")
		}
		if len(names) > 0 {
			for i, name := range names {
				names[i] = name[1:] // truncate the meta character.
			}
			nd.paramNames = names
		}
		nd.data = record.Value
	}
	return nil
}

func (da *DoubleArray) lookup(path string, params []string, idx int) (*node, []string) {
	if path == "" {
		return da.bc[idx], params
	}
	c, remaining := path[0], path[1:]
	if next, ok := da.check(c, idx); ok {
		if nd, params := da.lookup(remaining, params, next); nd != nil {
			return nd, params
		}
	}
	nd := da.bc[idx]
	if nd.paramTree != nil {
		i := nextSeparator(path, 0)
		remaining, params = path[i:], append(params, path[:i])
		if nd, params := nd.paramTree.lookup(remaining, params, 0); nd != nil {
			return nd, params
		}
	}
	if nd.isWildcard {
		return da.bc[idx], append(params, path)
	}
	return nil, nil
}

func (da *DoubleArray) build(routePaths []string, idx, depth int) error {
	base, siblings, err := da.arrange(routePaths, idx, depth)
	if err != nil {
		return err
	}
	for _, sib := range siblings {
		if !isMetaChar(sib.c) {
			da.setCheck(nextIndex(base, sib.c), idx)
		}
	}
	for _, sib := range siblings {
		switch sib.c {
		case paramCharacter:
			paths := routePaths[sib.start:sib.end]
			for i, path := range paths {
				paths[i] = path[nextSeparator(path, depth):]
			}
			rnd := da.bc[idx]
			if rnd.paramTree == nil {
				rnd.paramTree = New()
			}
			if err := rnd.paramTree.build(paths, 0, 0); err != nil {
				return err
			}
		case wildcardCharacter:
			da.bc[idx].isWildcard = true
		default:
			if err := da.build(routePaths[sib.start:sib.end], nextIndex(base, sib.c), depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// check returns next index of array of BASE/CHECK and whether the CHECK succeeded.
func (da *DoubleArray) check(c byte, i int) (next int, ok bool) {
	next = nextIndex(da.bc[i].base, c)
	return next, (da.bc[next] != nil && da.bc[next].check == i)
}

// setBase sets BASE.
func (da *DoubleArray) setBase(i, base int) {
	da.bc[i].base = base
}

// setCheck sets CHECK.
// If array of BASE/CHECK less than or equal to i, it will extend the array of BASE/CHECK.
func (da *DoubleArray) setCheck(i, check int) {
	if da.bc[i] == nil {
		da.bc[i] = newNode()
	}
	da.bc[i].check = check
}

// extendBaseCheckArray extends array of BASE/CHECK.
func (da *DoubleArray) extendBaseCheckArray() {
	da.bc = append(da.bc, make([]*node, blockSize)...)
}

// findBase returns good BASE.
func (da *DoubleArray) findBase(siblings []*sibling) (base int) {
	base = 1
	for count := 0; count < len(siblings); {
		for _, sib := range siblings {
			next := nextIndex(base, sib.c)
			if next >= len(da.bc) {
				da.extendBaseCheckArray()
			}
			if da.bc[next] != nil && (da.bc[next].check != -1 || da.bc[next].base != 0) {
				base++
				count = 0
				break
			}
			count++
		}
	}
	return base
}

func (da *DoubleArray) arrange(keys []string, idx, depth int) (base int, siblings []*sibling, err error) {
	siblings, err = makeSiblings(keys, depth)
	if err != nil {
		return -1, nil, err
	}
	if len(siblings) < 1 {
		return -1, nil, nil
	}
	base = da.findBase(siblings)
	da.setBase(idx, base)
	return base, siblings, err
}

// node represents a node of Double-Array.
type node struct {
	data  interface{}
	base  int
	check int

	// Tree of path parameter.
	paramTree *DoubleArray

	// Names of path parameters.
	paramNames []string

	// Whether the wildcard node.
	isWildcard bool
}

// newNode returns a new node.
func newNode() *node {
	return &node{base: 0, check: -1}
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

// isMetaChar returns whether the meta character.
func isMetaChar(c byte) bool {
	return c == paramCharacter || c == wildcardCharacter
}

// nextIndex returns next index of array of BASE/CHECK.
func nextIndex(base int, c byte) int {
	return base ^ int(c)
}

// nextSeparator returns an index of next separator in path.
func nextSeparator(path string, start int) int {
	for start < len(path) && path[start] != '/' && path[start] != '.' {
		start++
	}
	return start
}

// makeSiblings returns slice of sibling from string keys.
func makeSiblings(keys []string, depth int) (sib []*sibling, err error) {
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
			sib = append(sib, &sibling{start: i, c: c})
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
