package fmanager

import (
	"fmt"
	"io"
	"os"
	window "simpleSurance/MovingWindow"
	"time"
)

// FileManager is responsible for handling basic file operations.
type FileManager struct {
	f1In        *os.File
	f1Out       *os.File
	f2In        *os.File
	f2Out       *os.File
	LastChange1 int64
	LastChange2 int64
	FileName1   string
	FileName2   string
}

func (fm *FileManager) Close() error {
	err := fm.f1In.Close()
	err = fm.f2In.Close()
	err = fm.f1Out.Close()
	err = fm.f2Out.Close()

	return err
}

//We use 2 files, for backup
func (fm *FileManager) initFirstFile() error {
	var err error
	fm.f1In, err = os.OpenFile(fm.FileName1, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", fm.FileName1, err)
	}
	fm.f1Out, err = os.OpenFile(fm.FileName1, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", fm.FileName1, err)
	}
	return nil
}

func (fm *FileManager) initSecondFile() error {
	var err error
	fm.f2In, err = os.OpenFile(fm.FileName2, os.O_RDWR | os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", fm.FileName2, err)
	}
	fm.f2Out, err = os.OpenFile(fm.FileName2, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", fm.FileName2, err)
	}
	return nil
}

//Initialize backup files
func (fm *FileManager) Init() error{
	err := fm.initFirstFile()
	if err != nil{
		return err
	}
	err = fm.initSecondFile()
	return err
}

//
//updateWindow is responsible for reading time records from file
//and add it to moving window, returns last added valid time record
//
func updateWindow(f *os.File, mw *window.MovingWindow) int64{
	var t,change int64
	change = -1
	for {
		_ ,err := fmt.Fscanln(f, &t)
		if err != nil{
			break
		}
		change = t
		mw.Add(t)
	}
	return change
}

//
// UpdateFiles is responsible for adding new time record in files
// also it checks how old files are, clears all older than 1 minute files
// and updates last modification time for files
//
func (fm *FileManager) UpdateFiles(t int64) error{
	var err error
	if fm.LastChange1 < (t - time.Minute.Nanoseconds()) {
		err = fm.f1Out.Close()
		if err != nil{
			return fmt.Errorf("failed to close file %s: %s", fm.FileName1, err)
		}
		err = fm.f1In.Close()
		if err != nil{
			return fmt.Errorf("failed to close file %s: %s", fm.FileName1, err)
		}
		err = os.Remove(fm.FileName1)
		if err != nil{
			return fmt.Errorf("failed to remove file %s: %s", fm.FileName1, err)
		}
		err = fm.initFirstFile()
		if err != nil{
			return err
		}
		fm.LastChange1 = t
	}
	if fm.LastChange2 < t - time.Minute.Nanoseconds() {
		err = fm.f2Out.Close()
		if err != nil{
			return fmt.Errorf("failed to close file %s: %s", fm.FileName2, err)
		}
		err = fm.f2In.Close()
		if err != nil{
			return fmt.Errorf("failed to close file %s: %s", fm.FileName2, err)
		}
		err = os.Remove(fm.FileName2)
		if err != nil{
			return fmt.Errorf("failed to remove file %s: %s", fm.FileName2, err)
		}
		err = fm.initSecondFile()
		if err != nil{
			return err
		}
		fm.LastChange2 = t;
	}
	if fm.LastChange1 < fm.LastChange2 {
		_ ,err = fmt.Fprintln(fm.f2Out, t)
		if err != nil{
			return fmt.Errorf("failed to write in file %s: %s", fm.FileName2, err)
		}
		err = fm.f2Out.Sync()
		if err != nil{
			return err
		}
		fm.LastChange2 = t;
	} else {
		_ ,err = fmt.Fprintln(fm.f1Out, t)
		if err != nil{
			return fmt.Errorf("failed to write in file %s: %s", fm.FileName1, err)
		}
		err = fm.f1Out.Sync()
		if err != nil{
			return err
		}
		fm.LastChange1 = t;
	}

	return err
}

//
// Backup is responsible for backup records from files
// Records in files are always sorted
// Sets last change time for each files
//

func (fm *FileManager) Backup (mw *window.MovingWindow)  error{
	var a,b int64
	_ ,err := fmt.Fscanln(fm.f1In, &a)

	if (err != io.EOF) && (err != nil) {
		return fmt.Errorf("failed to read from file %s: %s", fm.FileName1, err)
	}
	_, err  = fmt.Fscanln(fm.f2In, &b)

	if (err != io.EOF) && (err != nil) {
		return fmt.Errorf("failed to read from file %s: %s", fm.FileName2, err)
	}
	//
	// read single time record from each file, compares them to each other
	// to define which one is older, notice that these 2 files can not have
	// intersections, since we always write in one file, when another file becomes
	// older than 1 minute, we clear it and change active file for recording
	// so all records in each file are sorted, and without intersection
	// and without mixing so we can get sorted time records just by concatenation
	//
	if a < b{
		if a != 0{
			mw.Add(a)
			fm.LastChange1 = a
		}
		val := updateWindow(fm.f1In, mw)
		if val > -1{
			fm.LastChange1 = val
		}

		if b != 0 {
			mw.Add(b)
			fm.LastChange2 = b
		}
		val = updateWindow(fm.f2In, mw)
		if val > -1{
			fm.LastChange2 = val
		}
	} else{
		if b != 0 {
			mw.Add(b)
			fm.LastChange2 = b
		}
		val := updateWindow(fm.f2In, mw)
		if val > -1 {
			fm.LastChange2 = val
		}
		if a != 0 {
			mw.Add(a)
			fm.LastChange1 = a
		}
		val = updateWindow(fm.f1In, mw)
		if val > -1 {
			fm.LastChange1 = val
		}
	}

	return nil
}