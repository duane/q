package main

import (
  "code.google.com/p/go-uuid/uuid"
  "fmt"
  "github.com/coreos/raft"
  "math/rand"
)

const commandPrefix = "q:"

func commandName(name string) string {
  return commandPrefix + name
}

type AddCommand struct {
  Message []byte
}

func (cmd *AddCommand) CommandName() string {
  return commandName("add")
}

func (cmd *AddCommand) Apply(raft *raft.Server) (interface{}, error) {
  server := raft.Context().(*Server)
  server.Queue.Enqueue(cmd.Message)
  return nil, nil
}

func main() {
  // random number between 19000-19999
  lower := 19000
  upper := 19999
  delta := upper - lower
  port := int(rand.Uint32()) % delta
  prefix := "/q"
  connectionString := fmt.Sprintf("localhost:%d", port)
  name := uuid.New()
  _, err := NewServer(name, prefix, connectionString)
  if err != nil {
    panic(err.Error())
  }
}
