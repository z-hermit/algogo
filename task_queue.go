package models

import (
	"fmt"
	"sync"
)

// OrderedBlock holds a block contents and its position on the file for queuing writer
type OrderedBlock struct {
	Priority int
	Key      interface{}
	Content  interface{}
	Position int
}

// A BlocksQueue is a min-heap of orderedBlocks.
type BlocksQueue struct {
	queue   []*OrderedBlock
	hashmap map[interface{}]*OrderedBlock
	lock    *sync.Mutex
	waiter  *sync.Cond
}

// New returns a pointer to an empty BlocksQueue ready to be used.
func NewQueue() *BlocksQueue {
	l := &sync.Mutex{}
	ret := &BlocksQueue{queue: []*OrderedBlock{{Priority: -1, Position: 0}}, hashmap: make(map[interface{}]*OrderedBlock), lock: l, waiter: sync.NewCond(l)}
	return ret
}

// len returns the number of elements in the queue
func (q BlocksQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.queue) - 1
}

func (q BlocksQueue) len() int {
	return len(q.queue) - 1
}

// Push add an element to the queue
func (q *BlocksQueue) Push(b *OrderedBlock) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if k, ok := q.hashmap[b.Key]; ok {
		k.Priority = b.Priority
		k.Content = b.Content
		q.up(k.Position)
		q.down(k.Position)
		return
	}

	q.hashmap[b.Key] = b
	b.Position = q.len() + 1
	q.queue = append(q.queue, b)
	q.up(b.Position)

	q.waiter.Signal()
}

func (q *BlocksQueue) up(position int) {
	idx := position
	parent := idx / 2

	for parent >= 1 && q.queue[idx].Priority > q.queue[parent].Priority {
		q.swap(idx, parent)
		idx = parent
		parent = idx / 2
	}
}

func (q *BlocksQueue) down(position int) {
	parent := position
	n := q.len()
	for {
		lChild := parent * 2
		if lChild > n || lChild < 1 {
			break
		}

		biggestChild := lChild

		rChild := lChild + 1
		if rChild <= n && q.queue[rChild].Priority > q.queue[lChild].Priority {
			biggestChild = rChild
		}

		if q.queue[parent].Priority > q.queue[biggestChild].Priority {
			break
		}
		q.swap(parent, biggestChild)
		parent = biggestChild
	}
}

func (q *BlocksQueue) swap(a, b int) {
	q.queue[a].Position = b
	q.queue[b].Position = a
	q.queue[a], q.queue[b] = q.queue[b], q.queue[a]
}

// Pop removes the lowest priority element from the queue and returns it
func (q *BlocksQueue) Pop() *OrderedBlock {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.len() == 0 {
		q.waiter.Wait()
	}

	ret := q.queue[1]
	n := q.len()
	q.queue[1] = q.queue[n]
	q.queue = q.queue[:n]
	// this keeps the reference to the original slice, but it's ok for my use case
	// I don't need to nil the left out element, as I will be checking for len() != 0
	delete(q.hashmap, ret.Key)
	if n != 1 {
		q.queue[1].Position = 1
		q.down(q.queue[1].Position)
	}

	return ret
}

func (q *BlocksQueue) Delete(key interface{}) *OrderedBlock {
	q.lock.Lock()
	defer q.lock.Unlock()
	if k, ok := q.hashmap[key]; ok {
		delete(q.hashmap, k)
		q.queue[k.Position] = q.queue[q.len()]
		q.queue = q.queue[:q.len()]
		if q.len()+1 != k.Position {
			q.queue[k.Position].Position = k.Position
			q.up(k.Position)
			q.down(k.Position)
		}
		return k
	}
	return nil
}

// Peek returns the lowest priority element from the queue
func (q BlocksQueue) Peek() *OrderedBlock {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.len() == 0 {
		return nil
	}
	return q.queue[1]
}

func (q BlocksQueue) Print() {
	fmt.Println("====")
	for _, v := range q.queue {
		fmt.Println(v)
	}
	fmt.Println("====")
}
