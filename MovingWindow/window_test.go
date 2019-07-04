package mwindow

import (
	"container/list"
	"testing"
	"time"
)

func TestMovingWindow_Request(t *testing.T) {
	mw := &MovingWindow{List: list.List{}}
	mw.Request(time.Now().UnixNano() - time.Minute.Nanoseconds())
	if mw.List.Len() != 0 {
		t.Errorf("List Size Not Match expected 0 was: %d", mw.List.Len())
	}
	mw.Request(time.Now().UnixNano())
	if mw.List.Len() != 1 {
		t.Errorf("List Size Not Match expected 1 was: %d", mw.List.Len())
	}
}

func TestMovingWindow_Clear(t *testing.T) {
	mw := &MovingWindow{List: list.List{}}
	mw.List.PushBack(time.Now().UnixNano())
	mw.Clear(time.Now().UnixNano())

	if mw.List.Len() != 1 {
		t.Errorf("List Size Not Match expected 1 was: %d", mw.List.Len())
	}

	mw.List.PushFront(time.Now().UnixNano() - time.Minute.Nanoseconds())
	mw.List.PushBack(time.Now().UnixNano())

	mw.Clear(time.Now().UnixNano())
	if mw.List.Len() != 2 {
		t.Errorf("List Size Not Match expected 2 was: %d", mw.List.Len())
	}
}

func isSorted (ls list.List) bool {
	var prev int64 = -1
	for e := ls.Front(); e != nil; e = e.Next() {
		v, _ := e.Value.(int64)
		if v < prev {
			return false
		}
		prev = v;
	}

	return true
}


//
//To imitate time elapsing we define now and after that we just add 1 second
//to nanoseconds to this now value
//
func TestMovingWindow_Moving_Mock(t *testing.T) {
	mw := &MovingWindow{List: list.List{}}
	now := time.Now().UnixNano()
	for i:=0; i<60; i++ {
		mw.List.PushBack(now)
		now += time.Second.Nanoseconds()
	}
	now += 500 * time.Millisecond.Nanoseconds()  //half second in nanoseconds

	for i:=59; i>10; i-- {
		mw.Clear(now)
		if mw.List.Len() != i {
			t.Errorf("List Size Not Match expected %d was: %d", i, mw.List.Len())
		}

		if !isSorted(mw.List) {
			t.Errorf("List is not sorted")
		}

		first, _ := mw.List.Front().Value.(int64)
		if now - first > time.Minute.Nanoseconds() {
			t.Errorf("Window contains old time record")
		}

		now += time.Second.Nanoseconds()
	}

	mw.List.PushBack(now)
	mw.Clear(now)
	if mw.List.Len() != 11 {
		t.Errorf("List Size Not Match expected %d was: %d", 11, mw.List.Len())
	}
}