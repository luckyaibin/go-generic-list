// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *List):
//
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
package list

import (
	"fmt"
	"strings"
)

// Element is an element of a linked list.
type Element[T any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[T]

	// The list to which this element belongs.
	list *List[T]

	// The value stored with this element.
	Value T
}

// Next returns the next list element or nil.
func (e *Element[T]) Next() *Element[T] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element[T]) Prev() *Element[T] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List[T any] struct {
	root Element[T] // sentinel list element, only &root, root.prev, and root.next are used
	len  int        // current list length excluding (this) sentinel element
}

// Init initializes or clears list l.
func (l *List[T]) Init() *List[T] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New[T any]() *List[T] { return new(List[T]).Init() }

func (l *List[T]) String() string {
	strbuilder := &strings.Builder{}
	if l.len == 0 {
		return ""
	}
	n := l.root.next
	for n != nil && n != &l.root {
		strbuilder.WriteString(fmt.Sprintf("%+v ", n.Value))
		n = n.next
	}
	return strbuilder.String()
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[T]) Len() int { return l.len }

// Front returns the first element of list l or nil if the list is empty.
func (l *List[T]) Front() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *List[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero List value.
func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

// remove removes e from its list, decrements l.len
func (l *List[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	l.len--
}

// move moves e to next to at.
func (l *List[T]) move(e, at *Element[T]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List[T]) Remove(e *Element[T]) T {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.remove(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	// see comment in List.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

func (l *List[T]) QuickSort(cmp func(a, b T) int) {
	l.lazyInit()
	first := l.Front()
	last := l.Back()
	_qsort(l, first, last, cmp)
}

func _qsort[T any](lst *List[T], left, right *Element[T], cmp func(a, b T) int) {
	if left == right {
		return
	}
	// LBoundary and RBoundary are boundaries before left and after right
	LBoundary := left.prev
	RBoundary := right.next
	pivotValue := left.Value
	var finalPivot *Element[T] = nil
	for {
		// Right moves first, and finally finds a value within the [left, right] interval that is<left or==left (only when it coincides with left)
		for right != left && cmp(right.Value, pivotValue) >= 0 {
			right = right.Prev()
		}
		// Then move left and finally find a value within the [left, right] interval that is>left or==left (only when it coincides with left)
		for left != right && cmp(left.Value, pivotValue) <= 0 {
			left = left.Next()
		}
		if left == right {
			break
		}
		isNeighbour := swap(left, right)
		left, right = right, left
		if isNeighbour { // When adjacent, there is no need for the next loop. In fact, there can be no break. The next time left and right are equal, there will also be a break
			break
		}
	}
	// After the loop ends, the position of left is the right boundary of all values less than or equal to pivotValue
	// Place pivotValue in the left position
	leftNext := left.next // Using the next position of the left, swapping and then taking the previous one is the final pivot position
	swap(LBoundary.next, left)
	finalPivot = leftNext.prev
	if LBoundary.next != finalPivot { // It may overlap to one point. Next time, recursion will occur, they cross over and loop to the wrong side and never ends.
		_qsort(lst, LBoundary.next, finalPivot.prev, cmp)
	}
	if finalPivot != RBoundary.prev { // Same as above
		_qsort(lst, finalPivot.next, RBoundary.prev, cmp)
	}
}
func swap[T any](b, d *Element[T]) (neighbor bool) {
	if b != d {
		if b.next == d || b.prev == d { // are neighours
			if b.next == d { // a <=> b <=> d <=> e
				b.prev.next = d
				d.next.prev = b
				d.prev = b.prev
				b.next = d.next
				d.next = b
				b.prev = d
			} else { // a <=> d <=> b <=> e
				d.prev.next = b
				b.next.prev = d
				b.prev = d.prev
				d.next = b.next
				b.next = d
				d.prev = b
			}
			neighbor = true
		} else { // a <=> b <=> c <=> d <=> e
			b.next.prev = d
			b.prev.next = d
			d.next.prev = b
			d.prev.next = b
			bPrevP := b.prev
			bNextP := b.next
			b.prev = d.prev
			b.next = d.next
			d.prev = bPrevP
			d.next = bNextP
		}
	}
	return
}
