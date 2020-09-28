package queue

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestNewQueue(t *testing.T) {
	qe := NewQueue(22)
	q := qe.(*queue)
	equalVal(t, 22, q.initCap, "length error")
	equalVal(t, 22, q.realCap, "length error")
	qe = NewQueue(0)
	equalVal(t, -1, q.initCap, "length error")
	equalVal(t, -1, q.realCap, "length error")
}

func TestQueue(t *testing.T) {
	qe := NewQueue(5)
	equalVal(t, 0, qe.Len(), "length error")
	equalVal(t, 5, qe.Cap(), "cap error")
	equalVal(t, true, qe.IsEmpty(), "isEmpty error")
	equalVal(t, false, qe.IsFull(), "isFull error")

	for i := 0; i < 5; i++ {
		err := qe.Push(i)
		equalVal(t, nil, err, "push errors")
	}

	equalVal(t, 5, qe.Len(), "length error")
	equalVal(t, 5, qe.Cap(), "cap error")
	equalVal(t, false, qe.IsEmpty(), "isEmpty error")
	equalVal(t, true, qe.IsFull(), "isFull error")

	err := qe.Push(5)
	equalVal(t, ErrIsFull, err, "push error")

	for i := 0; i < 5; i++ {
		val, err := qe.Pull()
		equalVal(t, nil, err, "pull error")
		equalVal(t, i, val, "pull error")
	}

	val, err := qe.Pull()
	equalVal(t, ErrIsEmpty, err, "pull error")
	equalVal(t, nil, val, "pull error")

	qe = NewQueue(0)
	equalVal(t, 0, qe.Len(), "length error")
	equalVal(t, -1, qe.Cap(), "cap error")
	equalVal(t, true, qe.IsEmpty(), "isEmpty error")
	equalVal(t, false, qe.IsFull(), "isFull error")

	for i := 0; i < 5; i++ {
		err := qe.Push(i)
		equalVal(t, nil, err, "push error")
	}

	equalVal(t, 5, qe.Len(), "length error")
	equalVal(t, qe.(*queue).realCap, qe.Cap(), "cap error")
	equalVal(t, false, qe.IsEmpty(), "isEmpty error")
	equalVal(t, false, qe.IsFull(), "isFull error")

	for i := 0; i < 5; i++ {
		val, err := qe.Pull()
		equalVal(t, nil, err, "pull error")
		equalVal(t, i, val, "pull error")
	}

	val, err = qe.Pull()
	equalVal(t, ErrIsEmpty, err, "pull error")
	equalVal(t, nil, val, "pull error")
}

func equalVal(t *testing.T, expect, val interface{}, errWithMsg string)  {
	if reflect.DeepEqual(expect, val) {
		return
	}
	_, file, line, _ := runtime.Caller(1)
	gopath := os.Getenv("GOPATH") + "/src/"
	if idx := strings.Index(file, gopath); idx != -1 {
		file = file[len(gopath):]
	}
	t.Error(fmt.Sprintf("[%s:%d] %s, val(%+v) should be equal to expect(%+v)", file, line, errWithMsg, val, expect))
}