package lock

import "testing"

func TestTryLock(t *testing.T) {
	m := TryLockMutex{}
	m.Lock()
	if m.TryLock() == true {
		t.Error("The state of acquiring the lock is wrong")
	}
	m.Unlock()
	if m.TryLock() == false {
		t.Error("The state of acquiring the lock is wrong")
		return
	}
	m.Unlock()
}
