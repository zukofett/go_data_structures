package main

type Data[T any] interface {
	any | List[T]
}

type Node[T any] struct {
	Data Data[T]
	next *Node[T]
}

func (n *Node[T]) Next() *Node[T] {
	if n != nil {
		return n.next
	}
	return nil
}

type List[T any] struct {
	head, tail *Node[T]
	length     int
}

func New[T any]() *List[T] {
	list := new(List[T])

	dummie := &Node[T]{
		next: nil,
		Data: list,
	}

	list.head = dummie
	list.tail = dummie

	return list
}

func (l *List[T]) Len() int {
	if l == nil {
		return 0
	}
	return l.length
}

func (l *List[T]) Begin() *Node[T] {
	if l == nil {
		return nil
	}
	return l.head
}

func (l *List[T]) End() *Node[T] {
	if l == nil {
		return nil
	}
	return l.tail
}

func (l *List[T]) Insert(data *T, at *Node[T]) *Node[T] {
	if l == nil {
		return nil
	}

	newNode := new(Node[T])

	if at.next == nil {
		list, ok := at.Data.(*List[T])
		if !ok {
			return nil
		}
		list.tail = newNode
	}

	newNode.next = at.next
	newNode.Data = at.Data

	at.next = newNode
	at.Data = data

    l.length++

	return at
}

func (l *List[T]) Remove(ele *Node[T]) *T {
	if l == nil {
		return nil
	}

	toRemove := ele.next

	if toRemove == nil {
		list, ok := ele.Data.(*List[T])
		if !ok {
			return nil
		}
		list.tail = ele
	}

	ele.next = toRemove.next
	ele.Data = toRemove.Data

	toRemove.next = nil
    l.length--

    data, ok := toRemove.Data.(*T)
    if !ok {
        return nil
    }
	return data
}


