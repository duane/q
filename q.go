package q

import (
  "github.com/coreos/raft"
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

func (cmd *AddCommand) Apply(server *raft.Server) (interface{}, error) {

}
