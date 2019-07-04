package fmanager

import (
	"container/list"
	"fmt"
	"os"
	mwindow "simpleSurance/MovingWindow"
	"testing"
	"time"
)

func TestFileManager_Backup(t *testing.T) {
	err := os.Remove("db1.txt")
	err = os.Remove("db2.txt")

	f1, err := os.OpenFile("db1.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	f2, err := os.OpenFile("db2.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)

	fileManager := &FileManager{LastChange1: 0, LastChange2:0,
		FileName1:"db1.txt", FileName2:"db2.txt"}
	err = fileManager.Init()

	if err!=nil {
		t.Errorf("Error During manager init: %s", err)
	}

	mw := &mwindow.MovingWindow{List: list.List{}}

	_, _ = fmt.Fprintln(f1, time.Now().UnixNano())
	err = fileManager.Backup(mw)
	if mw.List.Len() != 1 {
		t.Errorf("List Size Not Match expected 1 was: %d",mw.List.Len())
	}
	_, _ = fmt.Fprintln(f2, time.Now().UnixNano())
	err = fileManager.Backup(mw)
	if mw.List.Len() != 2 {
		t.Errorf("List Size Not Match expected 2 was: %d", mw.List.Len())
	}
	if err!=nil {
		t.Errorf("Error During backup: %s", err)
	}

	err = f1.Close()
	err = f2.Close()
	err = fileManager.Close()
	err = os.Remove("db1.txt")
	err = os.Remove("db2.txt")
}

func TestFileManager_UpdateFiles(t *testing.T) {
	err := os.Remove("db1.txt")
	err  = os.Remove("db2.txt")

	fileManager := &FileManager{LastChange1: 0, LastChange2:0,
		FileName1:"db1.txt", FileName2:"db2.txt"}
	err = fileManager.Init()

	if err!=nil {
		t.Errorf("Error During manager init: %s", err)
	}

	err = fileManager.UpdateFiles(time.Now().UnixNano())
	if err != nil {
		t.Errorf("Error During manager updateFiles: %s", err)
	}

	f1, err := os.OpenFile("db1.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	f2, err := os.OpenFile("db2.txt", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)

	var a,b int64
	_ ,err = fmt.Fscanln(f1, &a)

	_, err  = fmt.Fscanln(f2, &b)

	var c int
	if a!= 0{
		c++
	}
	if b != 0 {
		c++
	}
	if c != 1 {
		t.Errorf("Wrong result after manager updateFiles expected 1 was: %d", c)
	}

	err = f1.Close()
	err = f2.Close()
	err = fileManager.Close()
	err = os.Remove("db1.txt")
	err = os.Remove("db2.txt")
}

