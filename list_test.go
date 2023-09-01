package list

import (
	"math/rand"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	testCount := 1000
	for c := 0; c < testCount; c++ {
		maxCount := 10000
		lst := New[int]()
		for i := 0; i < maxCount; i++ {
			lst.PushBack(int(r.Intn(maxCount)))
		}
		// compare function make values increase
		lst.QuickSort(func(a, b int) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		})
		if lst.Len() != maxCount {
			t.Fail()
		}
		for curr := lst.Front(); ; {
			currV := curr.Value
			next := curr.Next()
			if next == nil {
				break
			}
			if currV > next.Value {
				t.Fatalf("invlid order %+v and %+v ", currV, next.Value)
			}
			curr = next
		}
	}
}
