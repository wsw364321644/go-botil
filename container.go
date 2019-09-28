package botil

import (
	"container/list"
	"sync"
)

type SafeList struct {
	l *list.List
	m *sync.RWMutex
}

func NewSafeList(mutex *sync.RWMutex) *SafeList {
	return &SafeList{l: list.New(),m:mutex}
}
func (sl *SafeList) SetLock(lock *sync.RWMutex) {
	sl.m=lock
}

func (sl *SafeList) PushBack(lock bool,v interface{}) {
	if v == nil {
		return
	}
	if lock{
		sl.m.Lock()
		defer sl.m.Unlock()
	}
	sl.l.PushBack(v)
}

func (sl *SafeList) Front(lock bool) *list.Element {
	if lock {
		sl.m.RLock()
		defer sl.m.RUnlock()
	}
	return sl.l.Front()
}

func (sl *SafeList) Remove(lock bool,e *list.Element) {
	if e == nil {
		return
	}
	if lock{
		sl.m.Lock()
		defer sl.m.Unlock()
	}
	sl.l.Remove(e)
}

func (sl *SafeList) Len() int {
	sl.m.RLock()
	defer sl.m.RUnlock()
	return sl.l.Len()
}

func (sl *SafeList) RemoveElem(lock bool,v interface{}) {
	if v == nil {
		return
	}
	if lock{
		sl.m.Lock()
		defer sl.m.Unlock()
	}

	for e := sl.l.Front(); e != nil;  {
		olde:=e
		e = e.Next()
		if olde.Value==v{
			sl.l.Remove(olde)
		}
	}
}