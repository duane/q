package main

import (
  "testing"
)

func TestNewQueue(t *testing.T) {
  q := NewQueue()
  if q == nil {
    t.Fail()
  }
}

func TestEnqueue(t *testing.T) {
  q := NewQueue()
  if q == nil {
    t.Fail()
  }

  q.Enqueue(4)
}

func TestDequeue(t *testing.T) {
  q := NewQueue()
  if q == nil {
    t.Fail()
  }

  q.Enqueue(4)
  val := q.Dequeue()
  if val.(int) != 4 {
    t.Fatalf("Expected 4, but got %+v", val)
  }
}

func TestQueueOrder(t *testing.T) {
  q := NewQueue()
  if q == nil {
    t.Fail()
  }

  q.Enqueue(4)
  q.Enqueue(5)
  if q.Dequeue().(int) != 4 || q.Dequeue().(int) != 5 {
    t.Fatal("Bad order!")
  }
}

func TestLen(t *testing.T) {
  q := NewQueue()
  if q.Len() != 0 {
    t.Fail()
  }

  q.Enqueue(4)
  if q.Len() != 1 {
    t.Fail()
  }

  q.Enqueue(5)
  if q.Len() != 2 {
    t.Fail()
  }

  q.Dequeue()
  if q.Len() != 1 {
    t.Fail()
  }

  q.Dequeue()
  if q.Len() != 0 {
    t.Fail()
  }
}

func TestMarshal(t *testing.T) {
  q := NewQueue()
  if q == nil {
    t.Fail()
  }

  q.Enqueue(4)
  q.Enqueue(5)
  data, err := q.Save()
  if err != nil {
    t.Fatal(err.Error())
    return
  }

  if data == nil {
    t.Fatal("Save should not return nil bytes when err is nil")
  }

  q = NewQueue()
  err = q.Recovery(data)
  if err != nil {
    t.Fatal(err.Error())
  }

  if q.Dequeue().(int) != 4 || q.Dequeue().(int) != 5 || q.Len() != 0 {
    t.Fatalf("Bad unmarshal!")
  }
}
