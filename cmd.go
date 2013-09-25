package main

import (
  "github.com/coreos/raft"
  "log"
)

const commandPrefix = "q:"

func commandName(name string) string {
  return commandPrefix + name
}

type EnqueueCommand struct {
  Message []byte
}

func (cmd *EnqueueCommand) CommandName() string {
  return commandName("enqueue")
}

func (cmd *EnqueueCommand) Apply(raft *raft.Server) (interface{}, error) {
  s.Queue.Enqueue(cmd.Message)
  return nil, nil
}

type DequeueCommand struct{}

func (cmd *DequeueCommand) CommandName() string {
  return commandName("dequeue")
}

func (cmd *DequeueCommand) Apply(raft *raft.Server) (interface{}, error) {
  payload := s.Queue.Dequeue()
  log.Printf("Dequeued in apply: %+v\n", payload)
  return payload, nil
}
