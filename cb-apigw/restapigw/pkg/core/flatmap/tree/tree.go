package tree

import "errors"

// ===== [ Constants and Variables ] =====

const (
	wildcard = "*"
)

var (
	errNoNilValuesAllowed = errors.New("no nil values allowed")
)

// ===== [ Types ] =====

// nodeAndPath - Tree 구조의 Node와 Path 정보를 표현하는 구조 형식
type nodeAndPath struct {
	n *node
	p []string
}

// edgeToMove - Tree 구조의 Edge 정보 이동을 위한 구조 형식
type edgeToMove struct {
	nodeAndPath
	e *edge
}

// Tree - Tree 구조의 Root Node 구조 형식
type Tree struct {
	root *node
}

// ===== [ Implementations ] =====

// collectMoveCandidates - 지정한 source path에 해당하는 이동 대상 node에 대한 nodeAndPath 정보를 구성하여 반환
func (t *Tree) collectMoveCandidates(srcPath []string, next []nodeAndPath) []nodeAndPath {
	acc := []nodeAndPath{}
	for _, step := range srcPath {
		if step == wildcard {
			for _, nap := range next {
				for _, e := range nap.n.edges {
					acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
				}
			}
		} else {
			for _, nap := range next {
				for _, e := range nap.n.edges {
					if step == e.label {
						acc = append(acc, nodeAndPath{n: e.n, p: append(nap.p, e.label)})
						break
					}
				}
			}
		}
		next, acc = acc, next[:0]
	}
	return next
}

// promoteEdges - 지정된 이동 대상 edge들을 지정된 경로로 이동
func (t *Tree) promoteEdges(edgesToMove []edgeToMove, destPath []string) {
	var l string
	lenDst := len(destPath)
	for _, n := range edgesToMove {
		parent := t.root
		for i, path := range destPath[:lenDst-1] {
			if path == wildcard {
				l = n.p[i]
			} else {
				l = path
			}

			found := false
			for _, e := range parent.edges {
				if e.label != l {
					continue
				}
				found = true
				parent = e.n
				break
			}

			if !found {
				break
			}
		}

		n.e.n.SetDepth(parent.depth + 1)
		n.e.label = destPath[lenDst-1]
		parent.edges = append(parent.edges, n.e)
	}
}

// embeddingEdges - 지정된 이동 대상 edge들을 지정된 경로로 삽입
func (t *Tree) embeddingEdges(edgesToMove []edgeToMove, destPath []string) {
	lenDst := len(destPath)
	for _, em := range edgesToMove {
		root := em.n
		for _, k := range destPath[:lenDst-1] {
			found := false
			for _, e := range root.edges {
				if e.label != k {
					continue
				}
				found = true
				root = e.n
				break
			}
			if found {
				continue
			}
			child := newNode(root.depth + 1)
			root.edges = append(root.edges, &edge{label: k, n: child})
			root = child
		}
		em.e.label = destPath[lenDst-1]
		em.e.n.SetDepth(root.depth + 1)
		root.edges = append(root.edges, em.e)
	}
}

// Add - Tree 구조에 지정한 Path에 지정한 값 설정
func (t *Tree) Add(path []string, v interface{}) {
	if v == nil {
		return
	}
	t.root.Add(path, v)
}

// Del - 지정된 Path 에 해당하는 모든 Node와 Edge들 삭제
func (t *Tree) Del(path []string) {
	t.root.Del(path...)
}

// Get - 지정된 Path에 해당하는 모든 Node와 Edge들 반환
func (t *Tree) Get(path []string) interface{} {
	return t.root.Get(path...)
}

// Move - 지정된 source Path에 해당하는 모든 Node와 Edge들을 지정된 destination Path로 이동
func (t *Tree) Move(srcPath, destPath []string) {
	next := []nodeAndPath{{n: t.root, p: []string{}}}

	prefixLen := len(srcPath)
	if prefixLen > 1 {
		next = t.collectMoveCandidates(srcPath[:prefixLen-1], next)
	}

	edgesToMove := []edgeToMove{}
	lenDst := len(destPath)

	isEdgeRelabel := prefixLen == lenDst

	for _, nap := range next {
		for i, e := range nap.n.edges {
			if e.label != srcPath[prefixLen-1] {
				continue
			}

			if isEdgeRelabel {
				e.label = destPath[prefixLen-1]
				break
			}

			edgesToMove = append(edgesToMove, edgeToMove{nodeAndPath: nap, e: e})

			copy(nap.n.edges[i:], nap.n.edges[i+1:])
			nap.n.edges[len(nap.n.edges)-1] = nil
			nap.n.edges = nap.n.edges[:len(nap.n.edges)-1]

			break
		}
	}

	if isEdgeRelabel {
		return
	}

	if prefixLen > lenDst {
		t.promoteEdges(edgesToMove, destPath)
		return
	}

	t.embeddingEdges(edgesToMove, destPath[prefixLen-1:])
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// New - 지정된 정보를 기준으로 새로운 Tree 생성
func New(v interface{}) (*Tree, error) {
	if v == nil {
		return nil, errNoNilValuesAllowed
	}

	tree := &Tree{
		root: &node{},
	}

	tree.Add([]string{}, v)
	return tree, nil
}
