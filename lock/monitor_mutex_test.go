package lock

import (
	"testing"
)

func TestMutex_IsLocked(t *testing.T) {
	m := MonitorMutex{}
	m.Lock()
	if m.Is(Locked) == false {
		t.Error("Locked, but the detected state is unlocked")
	}
	m.Unlock()
	if m.Is(Locked) == true {
		t.Error("Unlock, but the detected state is locked")
	}
}

/*
func TestMutex_Count(t *testing.T) {
	m := MonitorMutex{}
	for i := 0; i < 1000; i++ {
		go func() {
			m.Lock()
			time.Sleep(1 * time.Second)
			m.Unlock()
		}()
	}
	t.Logf("lock count: %d\n", m.Count())
	time.Sleep(5 * time.Second)
	t.Logf("lock count: %d", m.Count())
}
*/
