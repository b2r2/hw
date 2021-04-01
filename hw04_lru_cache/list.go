package hw04lrucache

// List ...
type List interface {
	Len() int
	Front() *Element
	Back() *Element
	PushFront(v interface{}) *Element
	PushBack(v interface{}) *Element
	Remove(i *Element)
	MoveToFront(i *Element)
}

// Element ...
type Element struct {
	Value interface{}
	next  *Element
	prev  *Element
}

// Prev ...
func (e Element) Prev() *Element {
	return e.prev
}

// Next ...
func (e Element) Next() *Element {
	return e.next
}

type list struct {
	head, tail *Element
	len        int
}

// List ...
func NewList() List {
	return new(list).init()
}

func (l *list) PushFront(v interface{}) *Element {
	if l.len == 0 {
		return l.insertFirstElement(&Element{Value: v})
	}
	return l.insertFrontElement(&Element{Value: v})
}

func (l *list) PushBack(v interface{}) *Element {
	if l.len == 0 {
		return l.insertFirstElement(&Element{Value: v})
	}
	return l.insertBackElement(&Element{Value: v})
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.head
}

func (l list) Back() *Element {
	if l.len == 0 {
		return nil
	}
	return l.tail
}

func (l *list) Remove(e *Element) {
	if e.next == nil && e.prev == nil {
		return
	}
	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
	if e == l.head {
		l.head = e.next
	}
	if e == l.tail {
		l.tail = e.prev
	}
	l.len--
}

func (l *list) MoveToFront(e *Element) {
	if e.next == nil && e.prev == nil {
		return
	}
	l.Remove(e)
	l.PushFront(e.Value)
}

func (l *list) init() List {
	l.len = 0
	l.head, l.tail = nil, nil
	return l
}
func (l *list) insertBackElement(e *Element) *Element {
	l.tail.next = e
	e.prev = l.tail
	l.tail = e
	l.len++
	return e
}

func (l *list) insertFrontElement(e *Element) *Element {
	l.head.prev = e
	e.next = l.head
	l.head = e
	l.len++
	return e
}

func (l *list) insertFirstElement(e *Element) *Element {
	l.head = e
	l.tail = e
	l.len++
	return e
}
