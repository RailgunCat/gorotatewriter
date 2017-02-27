package rotatewriter

import (
	"os"
	"sync"
	"time"
	"fmt"
	"strings"
	"errors"
)

type RotateWriter struct {
	lock             sync.Mutex
	filename         string
	fp               *os.File
	rotatingPeriod   time.Duration
	lastRotatingTime time.Time
}

func New(filename string, rotatingPeriod time.Duration) (w *RotateWriter, err error){

	if filename == "" {
		return nil, errors.New("invalid filename")
	}

	w = &RotateWriter{filename: filename, rotatingPeriod:rotatingPeriod, lastRotatingTime:time.Now()}
	err = w.Rotate()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	if w.rotationPeriodExpired(){
		err := w.Rotate()
		if err != nil {
			return -1, nil
		}
	}

	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.lastRotatingTime = time.Now()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		strDate := strings.Replace(time.Now().Format(time.RFC3339), ":", "-", -1) // for windows fs replace ":" to "-"
		err = os.Rename(w.filename, w.filename + "." + strDate)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	w.fp, err = os.Create(w.filename)
	return
}

func (w *RotateWriter) rotationPeriodExpired() bool {
	if w.lastRotatingTime.Add(w.rotatingPeriod).Before(time.Now()) {
		return true
	}
	return false
}
