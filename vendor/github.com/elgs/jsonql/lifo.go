package jsonql

import (
	"fmt"
)

// Lifo - last in first out stack
type Lifo struct {
	top  *Element
	size int
}

// Element - an item in the stack
type Element struct {
	value interface{}
	next  *Element
}

// Stack - a set of functions of the Stack
type Stack interface {
	Len() int
	Push(value interface{})
	Pop() (value interface{})
	Peep() (value interface{})
	Print()
}

// Len - gets the length of the stack.
func (s *Lifo) Len() int {
	return s.size
}

// Push - pushes the value into the stack.
func (s *Lifo) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

// Pop - pops the last value out of the stack.
func (s *Lifo) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

// Peep - gets the last value in the stack without popping it out.
func (s *Lifo) Peep() (value interface{}) {
	if s.size > 0 {
		value = s.top.value
		return
	}
	return nil
}

// Print - shows what's in the stack.
func (s *Lifo) Print() {
	tmp := s.top
	for i := 0; i < s.Len(); i++ {
		fmt.Print(tmp.value, ", ")
		tmp = tmp.next
	}
	fmt.Println()
}
