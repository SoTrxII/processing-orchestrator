package data_structures

import (
	"sync"
	"time"
)

type Dated interface {
	Timestamp() time.Time
}

// PriorityQueue represents a priority queue.
type PriorityQueue[I Dated] struct {
	mu      sync.Mutex
	items   []*I
	enqueue chan *I
	dequeue chan chan *I
	stop    chan struct{}
	wg      sync.WaitGroup
}

// Enqueue adds an item to the priority queue.
func (pq *PriorityQueue[I]) Enqueue(item *I) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.items = append(pq.items, item)
}

// Dequeue removes and returns the item with the earliest timestamp.
func (pq *PriorityQueue[I]) Dequeue() *I {
	responseChan := make(chan *I)
	pq.dequeue <- responseChan
	return <-responseChan
}

// processQueue is the main goroutine that manages the priority queue.
func (pq *PriorityQueue[I]) processQueue() {
	defer pq.wg.Done()

	for {
		select {
		case item := <-pq.enqueue:
			pq.mu.Lock()
			pq.items = append(pq.items, item)
			pq.mu.Unlock()
		case responseChan := <-pq.dequeue:
			pq.mu.Lock()
			if len(pq.items) == 0 {
				pq.mu.Unlock()
				responseChan <- nil
			} else {
				earliestTimestampIdx := 0
				for i, item := range pq.items {
					if (*item).Timestamp().Before((*pq.items[earliestTimestampIdx]).Timestamp()) {
						earliestTimestampIdx = i
					}
				}
				item := pq.items[earliestTimestampIdx]
				// Remove the item from the slice.
				pq.items = append(pq.items[:earliestTimestampIdx], pq.items[earliestTimestampIdx+1:]...)
				pq.mu.Unlock()
				responseChan <- item
			}
		case <-pq.stop:
			return
		}
	}
}
