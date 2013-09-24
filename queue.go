package q

import (
  "container/list"
)

type Atom []byte

type Queue struct {
  list *List
}

func (queue *Queue) Enqueue(x interface{}) {
  _ = queue.list.PushBack(x.(*Atom))
}

func (queue *Queue) Dequeue() interface{} {
  e := queue.list.Remove(queue.list.Front())
}

func (queue *Queue) Len() int {
  return len(queue.Heap)
}

func (queue *Queue) Iterate() (send <-chan *Atom, quit chan<- bool) {
  e := queue.list.Front()
  for e {

  }
}
