package tree

import (
	"cmp"
	"errors"
	//"slices"
)

var (
	ErrorKeyNotFound = errors.New("Key is not found")
)

type Color bool

const (
	RED   Color = true
	BLACK Color = false
)

type Node[K cmp.Ordered, V any] struct {
	Key    K
	Value  V
	Color  Color
	Left   *Node[K, V]
	Right  *Node[K, V]
	Parent *Node[K, V]
}

type RedBlackTree[K cmp.Ordered, V any] struct {
	Nil  *Node[K, V]
	Root *Node[K, V]
}

func NewTree[K cmp.Ordered, V any]() *RedBlackTree[K, V] {
	nilNode := &Node[K, V]{Color: BLACK}
	return &RedBlackTree[K, V]{
		Root: nilNode, Nil: nilNode,
	}
}
func (t *RedBlackTree[K, V]) leftRotate(x *Node[K, V]) {
	y := x.Right
	x.Right = y.Left
	if y.Left != t.Nil {
		y.Left.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == t.Nil {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	y.Left = x
	x.Parent = y
}

func (t *RedBlackTree[K, V]) rightRotate(x *Node[K, V]) {
	y := x.Left
	x.Left = y.Right
	if y.Right != t.Nil {
		y.Right.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == t.Nil {
		t.Root = y
	} else if x == x.Parent.Right {
		x.Parent.Right = y
	} else {
		x.Parent.Left = y
	}
	y.Right = x
	x.Parent = y
}

func (t *RedBlackTree[K, V]) transplant(x *Node[K, V], y *Node[K, V]) {
	if x.Parent == t.Nil {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	if y != t.Nil {
		y.Parent = x.Parent
	}
}

func (t *RedBlackTree[K, V]) insertFixup(z *Node[K, V]) {
	for z.Parent.Color == RED {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.leftRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.rightRotate(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y.Color == RED {
				z.Parent.Color = BLACK
				y.Color = BLACK
				z.Parent.Parent.Color = RED
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rightRotate(z)
				}
				z.Parent.Color = BLACK
				z.Parent.Parent.Color = RED
				t.leftRotate(z.Parent.Parent)
			}
		}
	}
	t.Root.Color = BLACK
}

func (t *RedBlackTree[K, V]) Search(key K) (*Node[K, V], bool) {
	current := t.Root
	for current != t.Nil && key != current.Key {
		if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	if current == t.Nil {
		return nil, false
	}
	return current, true
}

func (t *RedBlackTree[K, V]) minimum(x *Node[K, V]) *Node[K, V] {
	for x.Left != t.Nil {
		x = x.Left
	}
	return x
}

func (t *RedBlackTree[K, V]) removeFixup(x *Node[K, V]) {
	for x != t.Root && x.Color == BLACK {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == BLACK && w.Right.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Right.Color == BLACK {
					w.Left.Color = BLACK
					w.Color = RED
					t.rightRotate(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Right.Color = BLACK
				t.leftRotate(x.Parent)
				x = t.Root
			}
		} else {
			w := x.Parent.Left
			if w.Color == RED {
				w.Color = BLACK
				x.Parent.Color = RED
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Right.Color == BLACK && w.Left.Color == BLACK {
				w.Color = RED
				x = x.Parent
			} else {
				if w.Left.Color == BLACK {
					w.Right.Color = BLACK
					w.Color = RED
					t.leftRotate(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = BLACK
				w.Left.Color = BLACK
				t.rightRotate(x.Parent)
				x = t.Root
			}
		}
	}
	x.Color = BLACK
}

func (t *RedBlackTree[K, V]) height(node *Node[K, V]) int {
	if node == t.Nil {
		return 0
	}
	leftHeight := t.height(node.Left)
	rightHeight := t.height(node.Right)
	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

func (t *RedBlackTree[K, V]) validate(node *Node[K, V], currentBlackHeight int, targetBlackHeight *int) bool {
	if node == t.Nil {
		if *targetBlackHeight == -1 {
			*targetBlackHeight = currentBlackHeight
		} else if *targetBlackHeight != currentBlackHeight {
			return false
		}
		return true
	}
	if node.Color == RED {
		if node.Left.Color == RED || node.Right.Color == RED {
			return false
		}
	}
	nextBlackHeight := currentBlackHeight
	if node.Color == BLACK {
		nextBlackHeight++
	}
	return t.validate(node.Left, nextBlackHeight, targetBlackHeight) && t.validate(node.Right, nextBlackHeight, targetBlackHeight)
}

func (t *RedBlackTree[K, V]) InsertOrGet(key K) *Node[K, V] {
	current := t.Root
	parent := t.Nil
	for current != t.Nil {
		parent = current
		if key == current.Key {
			return current
		} else if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	z := &Node[K, V]{
		Key:    key,
		Color:  RED,
		Left:   t.Nil,
		Right:  t.Nil,
		Parent: parent,
	}
	if parent == t.Nil {
		t.Root = z
	} else if z.Key < parent.Key {
		parent.Left = z
	} else {
		parent.Right = z
	}
	t.insertFixup(z)
	return z
}

func (t *RedBlackTree[K, V]) Remove(key K) error {
	z := t.Nil
	current := t.Root
	for current != t.Nil {
		if key == current.Key {
			z = current
			break
		} else if key < current.Key {
			current = current.Left
		} else {
			current = current.Right
		}
	}
	if z == t.Nil {
		return ErrorKeyNotFound
	}
	y := z
	yOriginalColor := y.Color
	var x *Node[K, V]
	if z.Left == t.Nil {
		x = z.Right
		t.transplant(z, z.Right)
	} else if z.Right == t.Nil {
		x = z.Left
		t.transplant(z, z.Left)
	} else {
		y = t.minimum(z.Right)
		yOriginalColor = y.Color
		x = y.Right

		if y.Parent == z {
			x.Parent = y
		} else {
			t.transplant(y, y.Right)
			y.Right = z.Right
			y.Right.Parent = y
		}
		t.transplant(z, y)
		y.Left = z.Left
		y.Left.Parent = y
		y.Color = z.Color
	}
	if yOriginalColor == BLACK {
		t.removeFixup(x)
	}
	return nil
}

func (t *RedBlackTree[K, V]) InOrderWalk(fn func(key K, value V) bool) {
	var walk func(node *Node[K, V]) bool
	walk = func(node *Node[K, V]) bool {
		if node == t.Nil {
			return true
		}
		if !walk(node.Left) {
			return false
		}
		if !fn(node.Key, node.Value) {
			return false
		}
		return walk(node.Right)
	}
	walk(t.Root)
}

func (t *RedBlackTree[K, V]) Height() int {
	return t.height(t.Root)
}

func (t *RedBlackTree[K, V]) IsValid() bool {
	if t.Root != t.Nil && t.Root.Color != BLACK {
		return false
	}
	targetBlackHeight := -1
	return t.validate(t.Root, 0, &targetBlackHeight)
}
