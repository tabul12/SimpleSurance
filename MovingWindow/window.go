package mwindow

import (
	"container/list"
	"time"
)
//
//MovingWindow is responsible for saving new time records in
//in list, and also to remove old records, records in list are
//sorted (ascending order)
//

type MovingWindow struct {
	List list.List
}

func (w *MovingWindow) Add(t int64){
	w.List.PushBack(t)
	w.Clear(time.Now().UnixNano())
}

func (w *MovingWindow) Request(t int64) int{
	w.Add(t)
	l := w.List.Len()
	return l
}

func (mw *MovingWindow) Clear(now int64) {
	t := now - time.Minute.Nanoseconds()
	for e := mw.List.Front(); e != nil; {
		oldTR, _ := e.Value.(int64)
		if oldTR > t {
			//Since records in the list are sorted
			//next time records values will be bigger,
			//so it does not make a sense to continue iterating
			break
		}
		temp := e
		e = e.Next()
		mw.List.Remove(temp)
	}
}
