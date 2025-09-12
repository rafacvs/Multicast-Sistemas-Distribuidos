package pkg

import (
	"container/heap"
)

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

// Ordenação por: timestamp, ID remetente, sequência
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].DataTimestamp != pq[j].DataTimestamp {
		return pq[i].DataTimestamp < pq[j].DataTimestamp
	}

	if pq[i].ID.SenderID != pq[j].ID.SenderID {
		return pq[i].ID.SenderID < pq[j].ID.SenderID
	}

	return pq[i].ID.Seq < pq[j].ID.Seq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*QueueItem)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

type HoldBackQueue struct {
	queue PriorityQueue
}

func NewHoldBackQueue() *HoldBackQueue {
	q := &HoldBackQueue{
		queue: make(PriorityQueue, 0),
	}
	heap.Init(&q.queue)
	return q
}

func (q *HoldBackQueue) Push(item *QueueItem) {
	heap.Push(&q.queue, item)
}

func (q *HoldBackQueue) Pop() *QueueItem {
	if q.IsEmpty() {
		return nil
	}
	return heap.Pop(&q.queue).(*QueueItem)
}

func (q *HoldBackQueue) Peek() *QueueItem {
	if q.IsEmpty() {
		return nil
	}
	return q.queue[0]
}

func (q *HoldBackQueue) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *HoldBackQueue) Size() int {
	return len(q.queue)
}
