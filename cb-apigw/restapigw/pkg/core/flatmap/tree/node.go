package tree

import "strconv"

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// edge - Tree 구조의 Node 간 연결 정보 구조 형식
type edge struct {
	label string
	n     *node
}

// node - Tree 구조의 Node를 표현하는 구조 형식
type node struct {
	Value        interface{}
	isCollection bool
	edges        []*edge
	depth        int
}

// ===== [ Implementations ] =====

// flatten - 지정된 값을 Path가 없는 단순한 평면 값으로 설정
func (n *node) flatten(i interface{}) {
	switch v := i.(type) {
	case map[string]interface{}:
		n.isCollection = false
		if len(v) == 0 {
			n.Value = v
			break
		}

		for k, e := range v {
			n.Add([]string{k}, e)
		}
	case []interface{}:
		n.isCollection = true
		if len(v) == 0 {
			n.Value = v
			break
		}

		for i, e := range v {
			n.Add([]string{strconv.Itoa(i)}, e)
		}
	default:
		n.isCollection = false
		n.Value = v
	}
}

// Add - 지정된 경로에 지정된 값을 설정한 node 구성
func (n *node) Add(path []string, v interface{}) {
	if len(path) == 0 {
		n.flatten(v)
		return
	}

	for _, e := range n.edges {
		if e.label == path[0] {
			e.n.Add(path[1:], v)
			return
		}
	}

	child := newNode(n.depth + 1)
	n.edges = append(n.edges, &edge{label: path[0], n: child})
	child.Add(path[1:], v)
}

// Del - 지정된 Path에 해당하는 Node 삭제
func (n *node) Del(path ...string) {
	lenKs := len(path)

	if lenKs == 0 || n.IsLeaf() {
		return
	}

	if path[0] == wildcard {
		if lenKs > 1 {
			for _, e := range n.edges {
				e.n.Del(path[1:]...)
			}
			return
		}

		for i := range n.edges {
			n.edges[i] = nil
		}
		n.edges = n.edges[:0]
		return
	}

	for i, e := range n.edges {
		if e.label == path[0] {
			if lenKs == 1 {
				copy(n.edges[i:], n.edges[i+1:])
				n.edges[len(n.edges)-1] = nil
				n.edges = n.edges[:len(n.edges)-1]
				return
			}
			e.n.Del(path[1:]...)
			return
		}
	}
}

// Get - 지정된 Path에 해당하는 Node 정보 반환
func (n *node) Get(path ...string) interface{} {
	lenKs := len(path)
	lenEdges := len(n.edges)

	if lenEdges == 0 && lenKs > 0 {
		return nil
	}

	if lenKs == 0 {
		return n.expand()
	}

	if path[0] == wildcard {
		res := make([]interface{}, lenEdges)
		for i, e := range n.edges {
			res[i] = e.n.Get(path[1:]...)
		}
		return res
	}

	for _, e := range n.edges {
		if e.label == path[0] {
			return e.n.Get(path[1:]...)
		}
	}
	return nil
}

// IsLeaf - 해당 Node가 말단 Node인지 여부 반환
func (n *node) IsLeaf() bool {
	return len(n.edges) == 0
}

// expand - 해당 Node의 하위 노드들을 반환
func (n *node) expand() interface{} {
	children := len(n.edges)
	if children == 0 {
		return n.Value
	}

	if n.isCollection {
		res := make([]interface{}, children)
		for i, e := range n.edges {
			res[i] = e.n.Get()
		}

		return res
	}

	res := make(map[string]interface{}, children)
	for _, e := range n.edges {
		res[e.label] = e.n.Get()
	}
	return res
}

// SetDepth - 해당 Node의 Depth를 지정된 값으로 설정
func (n *node) SetDepth(d int) {
	n.depth = d
	for _, e := range n.edges {
		e.n.SetDepth(d + 1)
	}
}

// ===== [ Private Functions ] =====

// newNode - 지정한 Depth 정보를 가지는 Node 생성
func newNode(depth int) *node {
	return &node{
		edges: []*edge{},
		depth: depth,
	}
}

// ===== [ Public Functions ] =====
