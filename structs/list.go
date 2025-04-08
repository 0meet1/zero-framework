package structs

import (
	"container/list"
	"sync"
)

type ZeroLinked struct {
	list.List

	mutex sync.Mutex
}

func NewLinked() *ZeroLinked { return new(ZeroLinked).Init() }

func (zLinked *ZeroLinked) Init() *ZeroLinked {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.Init()
	return zLinked
}

func (zLinked *ZeroLinked) Len() int { return zLinked.List.Len() }

func (zLinked *ZeroLinked) Front() *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.Front()
}

func (zLinked *ZeroLinked) Back() *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.Back()
}

func (zLinked *ZeroLinked) Remove(e *list.Element) any {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.Remove(e)
}

func (zLinked *ZeroLinked) PushFront(v any) *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.PushFront(v)
}

func (zLinked *ZeroLinked) PushBack(v any) *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.PushBack(v)
}

func (zLinked *ZeroLinked) InsertBefore(v any, mark *list.Element) *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.InsertBefore(v, mark)
}

func (zLinked *ZeroLinked) InsertAfter(v any, mark *list.Element) *list.Element {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	return zLinked.List.InsertAfter(v, mark)
}

func (zLinked *ZeroLinked) MoveToFront(e *list.Element) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.MoveToFront(e)
}

func (zLinked *ZeroLinked) MoveToBack(e *list.Element) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.MoveToBack(e)
}

func (zLinked *ZeroLinked) MoveBefore(e, mark *list.Element) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.MoveBefore(e, mark)
}

func (zLinked *ZeroLinked) MoveAfter(e, mark *list.Element) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.MoveAfter(e, mark)
}

func (zLinked *ZeroLinked) PushBackList(other *list.List) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.PushBackList(other)
}

func (zLinked *ZeroLinked) PushFrontList(other *list.List) {
	zLinked.mutex.Lock()
	defer zLinked.mutex.Unlock()
	zLinked.List.PushFrontList(other)
}
