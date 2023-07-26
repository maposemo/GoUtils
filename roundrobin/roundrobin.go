package roundrobin

import "sync"

type RoundRobin struct {
	sync.RWMutex

	array      []any
	currentIdx int
}

func (rr *RoundRobin) Len() int {
	return len(rr.array)
}

func (rr *RoundRobin) Add(it any) {
	rr.Lock()
	defer rr.Unlock()

	rr.array = append(rr.array, it)
}

func (rr *RoundRobin) Remove(it any) {
	rr.Lock()
	defer rr.Unlock()

	for i := range rr.array {
		if rr.array[i] == it {
			rr.array = append(rr.array[:i], rr.array[i+1:]...)
		}
	}
}

func (rr *RoundRobin) Next() any {
	rr.RLock()
	defer rr.RUnlock()

	if len(rr.array) < 1 {
		return nil
	}

	next := rr.currentIdx + 1
	if next >= len(rr.array) {
		next = 0
	}

	rr.currentIdx = next
	return rr.array[next]
}
