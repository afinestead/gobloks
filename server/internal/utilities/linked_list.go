package utilities

type Node[T interface{}] struct {
	Value T
	Next  *Node[T]
}

type LinkedList[T interface{}] *Node[T]
