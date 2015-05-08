package store

import (
	"testing"
)

func TestPushFrontAndFront(t *testing.T) {
	cList := NewSafeLinkedList()

	key1 := "key1"
	key2 := "key2"
	cList.PushFront(key1)
	cList.PushFront(key2)

	front := cList.Front()
	size := cList.Len()
	if front.Value.(string) != key2 || size != 2 {
		t.Errorf("Error: TestPushFrontAndFront")
	}
}

func TestPopBackAndBack(t *testing.T) {
	cList := NewSafeLinkedList()

	key1 := "key1"
	key2 := "key2"
	cList.PushFront(key1)
	cList.PushFront(key2)

	back := cList.Back()
	size := cList.Len()
	if back.Value.(string) != key1 || size != 2 {
		t.Errorf("Error: TestPopBackAndBack")
	}

	back2 := cList.PopBack()
	size2 := cList.Len()
	if back2.Value.(string) != key1 || size2 != 1 {
		t.Errorf("Error: TestPopBackAndBack")
	}

}
