package ju

/* Node is a generic type implementing a node in a tree whose values are stored in type T. The values can be accessed through the struct field Val.
 */
type Node[T any] struct {
	parent,
	firstChild,
	left,
	right *Node[T]
	Val  T
	Flag map[any]any
}

// Parent returns parent of n
func (n *Node[T]) Parent() *Node[T] {

	return n.parent
}

// Left returns left sibling (if any) of n
func (n *Node[T]) Left() *Node[T] {

	return n.left
}

// Right returns left sibling (if any) of n
func (n *Node[T]) Right() *Node[T] {

	return n.right
}

// LastChild returns the last child if n. nil if n has no child.
func (n *Node[T]) LastChild() *Node[T] {

	if n.firstChild == nil {
		return nil
	}

	c := new(Node[T])

	for c = n.firstChild; c.right != nil; c = c.right {
	}

	return c
}

// Children returns the children of n in a slice.
func (n *Node[T]) Children() []*Node[T] {

	var resp []*Node[T]

	for c := n.firstChild; c != nil; c = c.right {
		resp = append(resp, c)
	}

	return resp
}

// IsFirstSibling returns true iff. n is left-most sibling.
func (n *Node[T]) IsFirstSibling() bool {

	return n.left == nil
}

// IsLastSibling returns true iff. n is right-most sibling.
func (n *Node[T]) IsLastSibling() bool {

	return n.right == nil
}

// LastSibling returns the right-most sibling of n.
func (n *Node[T]) LastSibling() *Node[T] {

	s := new(Node[T])

	for s = n; s.right != nil; s = s.right {
	}
	return s
}

// FirstSibling returns the left-nmost sibling of n.
func (n *Node[T]) FirstSibling() *Node[T] {

	s := new(Node[T])

	for s = n; s.left != nil; s = s.left {
	}

	return s
}

// ChildCount returns the number of child nodes of n.
func (n *Node[T]) ChildCount() int {

	if n.firstChild == nil {
		return 0
	}

	i := 1

	for c := n.firstChild; c != nil; c = c.right {
		i++
	}

	return i
}

// SiblingCount returns the number of siblings of n.
func (n *Node[T]) SiblingCount() int {

	c := 0

	s := new(Node[T])

	for s = n; s != nil; s = s.left {
		c++
	}

	for s = n.right; s != nil; s = s.right {
		c++
	}

	return c
}

// Child returns the i-th child of n. WE START COUNTING AT CHILD 0 (not 1).
func (n *Node[T]) Child(i int) *Node[T] {

	if n.firstChild == nil {
		return nil
	}

	for k, c := 0, n.firstChild; c != nil; k, c = k+1, c.right {
		if k == i {
			return c
		}
	}

	return nil
}

// Root returns the root node of the tree to which n belongs.
func (n *Node[T]) Root() *Node[T] {

	e := new(Node[T])

	for e = n; e.parent != nil; e = e.parent {
	}

	return e
}

// AddChild adds n2 as the last child of n.
func (n *Node[T]) AddChild(n2 *Node[T]) {

	if n.firstChild == nil {
		n.firstChild = n2
	} else {
		n.LastChild().right = n2
	}

	n2.parent = n
}

// InsertRight inserts n2 to the right of n. n2's parent will be the
// same as as n's.
func (n *Node[T]) InsertRight(n2 *Node[T]) {

	nr := n.right

	n.right = n2
	n2.left = n
	n2.right = nr
	if nr != nil {
		nr.left = n2
	}

	n2.parent = n.parent
}

// InsertLeft inserts n2 to the left of n. n2's parent will be the same as n's.
func (n *Node[T]) InsertLeft(n2 *Node[T]) {

	nl := n.left

	n.left = n2
	n2.right = n
	n2.left = nl
	if nl != nil {
		nl.right = n2
	}

	n2.parent = n.parent

	if n2.parent != nil && n2.left == nil {
		n.parent.firstChild = n2
	}
}

// Remove n from tree. The siblings to the left and right of n are immediate siblings
// afterwards. n will have no siblings and no parent but keeps its children so you get a detached tree rooted in n.
func (n *Node[T]) Remove() {

	nl := n.left
	nr := n.right

	if nl != nil {
		nl.right = nr
	}

	if nr != nil {
		nr.left = nl
	}

	n.right = nil
	n.left = nil
	n.parent = nil
}

// Replace n with n2. n will be the root of a detached tree afterwards.
func (n *Node[T]) ReplaceWith(n2 *Node[T]) {

	n2.parent = n.parent
	n2.left = n.left
	n2.right = n.right

	if n2.left != nil {
		n2.left.right = n2
	}

	if n2.right != nil {
		n2.right.left = n2
	}

	n.parent = nil
	n.left = nil
	n.right = nil
}

// RemoveChildren removes the children of n. All the child nodes of n will have no parent afterwards.
func (n *Node[T]) RemoveChildren() {

	if n.firstChild == nil {
		return
	}

	for c := n.firstChild; c != nil; c = c.right {
		c.parent = nil
	}

	n.firstChild = nil
}

// SetFlag sets a flag with key k and value v.
func (n *Node[T]) SetFlag(k any, v ...any) {

	n.ensureFlagInitialized()

	switch len(v) {

	case 0:
		n.Flag[k] = nil

	case 1:
		n.Flag[k] = v[0]

	default:
		n.Flag[k] = v
	}
}

// HasFlag returns true if a flag with key k has been set.
func (n *Node[T]) HasFlag(k any) bool {

	n.ensureFlagInitialized()

	_, ok := n.Flag[k]
	return ok
}

// RemoveFlag removes the flag with key k.
func (n *Node[T]) RmFlag(k any) {

	n.ensureFlagInitialized()

	delete(n.Flag, k)
}

func (n *Node[T]) ensureFlagInitialized() {

	if n.Flag == nil {
		n.Flag = make(map[any]any)
	}
}

// Walk the tree in the standard depth-first, left-to-right order while executing the supplied functions. ingressFunc is executed when entering a node, egressFunc when exiting a node, pivotFunc when switching from one child to the next.
func (n *Node[T]) Walk(ingressFunc, pivotFunc, egressFunc func(*Node[T])) {

	var walk func(*Node[T])

	walk = func(e *Node[T]) {

		if e == nil {
			return
		}

		ingressFunc(e)

		for i, c := range e.Children() {

			walk(c)

			if i < len(e.Children())-1 {
				pivotFunc(e)
			}

		}

		egressFunc(e)

	}

	walk(n)
}

// Linearize the tree and return it as a slice of nodes. The order is standard depth-first, left-to-right.
func (n *Node[T]) Linearize() []*Node[T] {

	var resp []*Node[T]

	ingressFunc := func(e *Node[T]) {

		if e == nil {
			return
		}

		resp = append(resp, e)

	}

	n.Walk(ingressFunc, doNothing, doNothing)

	return resp

}

// Linearize the tree and return it as a slice but in the reverse-Polish style order.
func (n *Node[T]) LinearizeReverse() []*Node[T] {

	var resp []*Node[T]

	egressFunc := func(e *Node[T]) {

		if e == nil {
			return
		}

		resp = append(resp, e)

	}

	n.Walk(doNothing, doNothing, egressFunc)

	return resp

}

// Height returns the height of the tree rooted at n.
func (n *Node[T]) Height() int {

	h := -1
	ch := -1

	ingressFunc := func(e *Node[T]) {
		ch++
		if ch > h {
			h = ch
		}
	}

	egressFunc := func(e *Node[T]) {
		ch--
	}

	n.Walk(ingressFunc, doNothing, egressFunc)

	return h
}

// Depth returns the depth of the node (i.e., distance to root).
func (n *Node[T]) Depth() int {

	d := -1

	for e := n; e != nil; e = e.parent {
		d++
	}

	return d
}

func doNothing[T any](e *Node[T]) {
	return
}

var customFuncs struct {
	stringer func(any) string
}

func init() {
}

// Check is a convenience wrapper for the function test.
func (n *Node[T]) Check(test func(*Node[T]) bool) bool {

	return test(n)
}

// FindFunc is a convenience wrapper for the find-algorithm implemented by the function f.
func (n *Node[T]) FindFunc(f func(*Node[T]) []*Node[T]) []*Node[T] {

	return f(n)
}

// FindMatch returns all nodes that match the condition supplied by matchFunc.
func (n *Node[T]) FindMatch(matchFunc func(a *Node[T]) bool) []*Node[T] {

	var resp []*Node[T]

	ingressFunc := func(e *Node[T]) {
		if matchFunc(e) == true {
			resp = append(resp, e)
		}
	}

	n.Walk(ingressFunc, doNothing, doNothing)

	return resp
}

// SetStringFunc sets the function for determining the return value of the String method of n.
func SetStringFunc[T any](f func(*Node[T]) string) {

	customFuncs.stringer = func(n any) string {
		m, _ := n.(*Node[T])
		return f(m)
	}
}

// String esures n is a stringer. NOTE: you must first set the stringerFunction using [SetStringerFunc]
func (n *Node[T]) String() string {

	return customFuncs.stringer(n)
}

// WalkLevelOrder performs a breadth-first, left-to-right walk which executing the two supplied functions. ingressFunc is executed when processing a node, descendFunc is executed when entering a level (so it's the first function to be executed in any walk of this style.
func (n *Node[T]) WalkLevelOrder(ingressFunc, descendFunc func(n *Node[T])) {

	var children []*Node[T]
	var tmp []*Node[T]

	for children = append(children, n); len(children) != 0; {

		descendFunc(children[0])

		tmp = nil
		for _, e := range children {
			ingressFunc(e)
			tmp = append(tmp, e.Children()...)
		}

		children = nil
		children = append(children, tmp...)
	}
}
