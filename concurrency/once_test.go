package concurrency

import (
	"errors"
	"testing"
)

func TestOnce(t *testing.T) {
	count := 1
	f := func() error {
		if count > 0 {
			count--
			return nil
		}
		return errors.New("count value should more than 0")
	}

	once := Once{}
	if err := once.Do(f); err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("'do' function does not work")
	}
	if !once.Done() {
		t.Error("incorrect status for do")
	}

	if err := once.Do(f); err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("do function repeated execution")
	}
	if !once.Done() {
		t.Error("incorrect status for do")
	}

}