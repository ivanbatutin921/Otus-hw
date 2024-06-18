package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	List  // Remove me after realization.
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	if l.front == nil {
		l.front = newNode
		l.back = newNode
	} else {
		newNode.Next = l.front
		l.front.Prev = newNode
		l.front = newNode
	}
	l.len++
	return newNode
}

func (l *list) PushBack(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	if l.back == nil {
		l.front = newNode
		l.back = newNode
	} else {
		newNode.Prev = l.back
		l.back.Next = newNode
		l.back = newNode
	}
	l.len++
	return newNode
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.back = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}
