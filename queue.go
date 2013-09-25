package main

import (
  "container/list"
)

type Queue struct {
  List *list.List
}

func NewQueue() (queue *Queue) {
  queue = &Queue{List: list.New()}
  queue.List = queue.List.Init()
  return
}

func (queue *Queue) Enqueue(x interface{}) {
  _ = queue.List.PushBack(&x)
}

func (queue *Queue) Dequeue() interface{} {
  return *queue.List.Remove(queue.List.Front()).(*interface{})
}

func (queue *Queue) Len() int {
  return queue.List.Len()
}
