// A URL router implemented by Double-Array Trie.
package doublearray

import (
	"fmt"
	"sort"
	"strings"

	"github.com/naoina/kocha-urlrouter"
)

const (
	// Block size of array of BASE/CHECK of Double-Array.
	blockSize = 256
)

// baseCheck represents a BASE/CHECK node.
type baseCheck struct {
	base      int
	check     int
	hasParams bool
}

// DoubleArray represents a URLRouter by Double-Array.
type DoubleArray struct {
	static *doubleArray
	param  *doubleArray
}

// New returns a new DoubleArray.
func New() *DoubleArray {
	return &DoubleArray{
		static: newDoubleArray(blockSize),
		param:  newDoubleArray(blockSize),
	}
}

type doubleArray struct {
	bc   []baseCheck
	node map[int]*node
}

func newDoubleArray(size int) *doubleArray {
	return &doubleArray{
		bc:   newBaseCheckArray(size),
		node: make(map[int]*node),
	}
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
	if idx, found := da.static.lookupStatic(path); found {
		return da.static.node[idx].data, nil
	}
	nodes, idx, values := da.param.lookupParam(path, nil)
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
func (da *DoubleArray) Build(records []urlrouter.Record) error {
	statics, params := makeRecords(records)
	if err := da.static.build(statics, 0, 0); err != nil {
		return err
	}
	if err := da.param.build(params, 0, 0); err != nil {
		return err
	}
	return nil
}

func (da *doubleArray) lookupStatic(path string) (idx int, found bool) {
	for i := 0; i < len(path); i++ {
		next := nextIndex(da.bc[idx].base, path[i])
		if da.bc[next].check != idx {
			return -1, false
		}
		idx = next
	}
	return idx, true
}

func (da *doubleArray) lookupParam(path string, params []string) (map[int]*node, int, []string) {
	idx := 0
	var indexes []int64
	for i := 0; i < len(path); i++ {
		next := nextIndex(da.bc[idx].base, path[i])
		if da.bc[next].check != idx {
			goto PARAMED_ROUTE
		}
		idx = next
		if da.bc[idx].hasParams {
			indexes = append(indexes, int64(((i+1)&0xffffffff)<<32)|int64(idx&0xffffffff))
		}
	}
	return da.node, idx, params
PARAMED_ROUTE:
	for i := len(indexes) - 1; i >= 0; i-- {
		curIdx, idx := int((indexes[i]>>32)&0xffffffff), int(indexes[i]&0xffffffff)
		nd := da.node[idx]
		if nd.paramTree != nil {
			i := urlrouter.NextSeparator(path, curIdx)
			remaining, params := path[i:], append(params, path[curIdx:i])
			if nodes, idx, params := nd.paramTree.lookupParam(remaining, params); nodes != nil {
				return nodes, idx, params
			}
		}
		if nd.wildcardTree != nil {
			return nd.wildcardTree.node, 0, append(params, path[curIdx:])
		}
	}
	return nil, -1, nil
}

func (da *doubleArray) build(srcs []*Record, idx, depth int) error {
	base, siblings, leaf, err := da.arrange(srcs, idx, depth)
	if err != nil {
		return err
	}
	if leaf != nil {
		nd, err := makeNode(leaf)
		if err != nil {
			return err
		}
		da.node[idx] = nd
	}
	for _, sib := range siblings {
		if !urlrouter.IsMetaChar(sib.c) {
			da.setCheck(nextIndex(base, sib.c), idx)
		}
	}
	for _, sib := range siblings {
		switch records := srcs[sib.start:sib.end]; sib.c {
		case urlrouter.ParamCharacter:
			for _, record := range records {
				next := urlrouter.NextSeparator(record.Key, depth)
				name := record.Key[depth+1 : next]
				record.paramNames = append(record.paramNames, name)
				record.Key = record.Key[next:]
			}
			if da.node[idx] == nil {
				da.node[idx] = &node{}
			}
			da.node[idx].paramTree = newDoubleArray(blockSize)
			da.bc[idx].hasParams = true
			if err := da.node[idx].paramTree.build(records, 0, 0); err != nil {
				return err
			}
		case urlrouter.WildcardCharacter:
			if da.node[idx] == nil {
				da.node[idx] = &node{}
			}
			record := records[0]
			name := record.Key[depth+1:]
			record.paramNames = append(record.paramNames, name)
			da.node[idx].wildcardTree = newDoubleArray(0)
			nd, err := makeNode(record)
			if err != nil {
				return err
			}
			da.node[idx].wildcardTree.node[0] = nd
			da.bc[idx].hasParams = true
		default:
			if err := da.build(records, nextIndex(base, sib.c), depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// setBase sets BASE.
func (da *doubleArray) setBase(i, base int) {
	da.bc[i].base = base
}

// setCheck sets CHECK.
func (da *doubleArray) setCheck(i, check int) {
	da.bc[i].check = check
}

// extendBaseCheckArray extends array of BASE/CHECK.
func (da *doubleArray) extendBaseCheckArray() {
	da.bc = append(da.bc, newBaseCheckArray(blockSize)...)
}

// findEmptyIndex returns an index of unused BASE/CHECK node.
func (da *doubleArray) findEmptyIndex(start int) int {
	i := start
	for ; i < len(da.bc); i++ {
		if da.bc[i].base == 0 && da.bc[i].check == -1 {
			break
		}
	}
	return i
}

// findBase returns good BASE.
func (da *doubleArray) findBase(siblings []sibling, start int) (base int) {
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

func (da *doubleArray) arrange(records []*Record, idx, depth int) (base int, siblings []sibling, leaf *Record, err error) {
	siblings, leaf, err = makeSiblings(records, depth)
	if err != nil {
		return -1, nil, nil, err
	}
	if len(siblings) < 1 {
		return -1, nil, leaf, nil
	}
	base = da.findBase(siblings, idx)
	da.setBase(idx, base)
	return base, siblings, leaf, err
}

// node represents a node of Double-Array.
type node struct {
	data interface{}

	// Tree of path parameter.
	paramTree *doubleArray

	// Tree of wildcard path parameter.
	wildcardTree *doubleArray

	// Names of path parameters.
	paramNames []string
}

// makeNode returns a new node from record.
func makeNode(record *Record) (*node, error) {
	dups := make(map[string]bool)
	for _, name := range record.paramNames {
		if dups[name] {
			return nil, fmt.Errorf("path parameter `%v` is duplicated in the key '%v'", name, record.Key)
		}
		dups[name] = true
	}
	return &node{data: record.Value, paramNames: record.paramNames}, nil
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

// makeSiblings returns slice of sibling.
func makeSiblings(records []*Record, depth int) (sib []sibling, leaf *Record, err error) {
	var (
		pc byte
		n  int
	)
	for i, record := range records {
		if len(record.Key) == depth {
			leaf = record
			continue
		}
		c := record.Key[depth]
		switch {
		case pc < c:
			sib = append(sib, sibling{start: i, c: c})
		case pc == c:
			continue
		default:
			return nil, nil, fmt.Errorf("BUG: routing table hasn't been sorted")
		}
		if n > 0 {
			sib[n-1].end = i
		}
		pc = c
		n++
	}
	if n == 0 {
		return nil, leaf, nil
	}
	sib[n-1].end = len(records)
	return sib, leaf, nil
}

// Record represents a record that use to build the Double-Array.
type Record struct {
	urlrouter.Record
	paramNames []string
}

// RecordSlice represents a slice of Record for sort and implements the sort.Interface.
type RecordSlice []*Record

// makeRecords returns the records that use to build Double-Arrays.
func makeRecords(srcs []urlrouter.Record) (statics, params []*Record) {
	spChars := string([]byte{urlrouter.ParamCharacter, urlrouter.WildcardCharacter})
	for _, record := range srcs {
		if strings.ContainsAny(record.Key, spChars) {
			params = append(params, &Record{Record: record})
		} else {
			statics = append(statics, &Record{Record: record})
		}
	}
	sort.Sort(RecordSlice(statics))
	sort.Sort(RecordSlice(params))
	return statics, params
}

// Len implements the sort.Interface.Len.
func (rs RecordSlice) Len() int {
	return len(rs)
}

// Less implements the sort.Interface.Less.
func (rs RecordSlice) Less(i, j int) bool {
	return rs[i].Key < rs[j].Key
}

// Swap implements the sort.Interface.Swap.
func (rs RecordSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
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
