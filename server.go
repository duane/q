package main

import (
  "github.com/coreos/raft"
  "net/http"
  "path"
)

type Server struct {
  Name   string
  Prefix string
  Raft   *raft.Server
  Queue  *Queue
}

func NewServer(name string, prefix string, connectionString string) (server *Server, err error) {
  server = &Server{Name: name, Prefix: prefix, Queue: NewQueue()}

  raft_path := path.Join(prefix, "raft")
  transporter := raft.NewHTTPTransporter(raft_path)
  server.Raft, err = raft.NewServer(name, raft_path, transporter, server.Queue, server, connectionString)
  if err != nil {
    return
  }

  transporter.Install(server.Raft, http.DefaultServeMux)

  return
}
